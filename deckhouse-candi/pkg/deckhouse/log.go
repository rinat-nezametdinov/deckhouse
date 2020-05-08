package deckhouse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flant/logboek"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"flant/deckhouse-candi/pkg/kube"
)

type logLine struct {
	Module  string `json:"module,omitempty"`
	Level   string `json:"level,omitempty"`
	Output  string `json:"output,omitempty"`
	Message string `json:"msg,omitempty"`
}

func PrintDeckhouseLogs(client *kube.KubernetesClient, stopChan *chan struct{}) error {
	pods, err := client.CoreV1().Pods("d8-system").List(metav1.ListOptions{LabelSelector: "app=deckhouse"})
	if err != nil {
		return err
	}

	if len(pods.Items) != 1 {
		return fmt.Errorf("only one deckhouse pod should exist: %v", pods.Items)
	}

	logOptions := corev1.PodLogOptions{Container: "deckhouse", TailLines: int64Pointer(5)}

	timer := time.NewTicker(3 * time.Second)
	defer timer.Stop()

	printDots := true

	stopPrintDots := func() {
		if printDots {
			logboek.LogInfoF("\n")
			printDots = false
		}
	}

	for {
		select {
		case <-*stopChan:
			return nil
		case <-timer.C:
			if printDots {
				logboek.LogInfoF(".")
			}

			request := client.CoreV1().Pods("d8-system").GetLogs(pods.Items[0].Name, &logOptions)
			result, err := request.DoRaw()
			if err != nil {
				return err
			}

			currentTime := metav1.NewTime(time.Now())
			logOptions = corev1.PodLogOptions{Container: "deckhouse", SinceTime: &currentTime}

			reader := bufio.NewReader(bytes.NewReader(result))
			for {
				l, _, err := reader.ReadLine()
				if err != nil {
					break
				}
				var line logLine
				if err := json.Unmarshal(l, &line); err != nil {
					return err
				}

				if line.Level == "error" || line.Output == "stderr" {
					stopPrintDots()
					logboek.LogWarnLn(line.Message)
					continue
				}

				if line.Message == "Module run success" {
					stopPrintDots()
					logboek.LogInfoF("Module %q run successfully\n", line.Module)
					continue
				}
			}
		}
	}
}

func int64Pointer(i int) *int64 {
	r := int64(i)
	return &r
}