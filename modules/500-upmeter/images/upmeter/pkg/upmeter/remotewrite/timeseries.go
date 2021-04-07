package remotewrite

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/prometheus/prompb"

	"upmeter/pkg/check"
)

func convEpisodes2Timeseries(timeslot time.Time, episodes []*check.DowntimeEpisode, commonLabels []*prompb.Label) []*prompb.TimeSeries {
	tss := make([]*prompb.TimeSeries, 0)

	for _, ep := range episodes {
		labels := []*prompb.Label{}
		labels = append(labels, episodeLabels(ep)...)
		labels = append(labels, commonLabels...)

		nodata := ep.NoDataSeconds
		fail := ep.FailSeconds
		unknown := ep.UnknownSeconds
		success := ep.SuccessSeconds

		tss = append(tss,
			statusTimeseries(timeslot, success, withLabel(labels, &prompb.Label{Name: "status", Value: "up"})),
			statusTimeseries(timeslot, fail, withLabel(labels, &prompb.Label{Name: "status", Value: "down"})),
			statusTimeseries(timeslot, unknown, withLabel(labels, &prompb.Label{Name: "status", Value: "unknown"})),
			statusTimeseries(timeslot, nodata, withLabel(labels, &prompb.Label{Name: "status", Value: "nodata"})),
		)
	}
	return tss
}

func statusTimeseries(timeslot time.Time, value int64, labels []*prompb.Label) *prompb.TimeSeries {
	return &prompb.TimeSeries{
		Labels: labels,
		Samples: []prompb.Sample{
			{
				Timestamp: timeslot.Unix() * 1e3, // milliseconds
				Value:     float64(value * 1e3),  // milliseconds
			},
		},
	}
}

func withLabel(originalLabels []*prompb.Label, statusLabel *prompb.Label) []*prompb.Label {
	labels := make([]*prompb.Label, len(originalLabels), len(originalLabels)+1)
	copy(labels, originalLabels)
	labels = append(labels, statusLabel)

	return labels
}

func episodeLabels(ep *check.DowntimeEpisode) []*prompb.Label {
	return []*prompb.Label{
		{
			Name:  "__name__",
			Value: "statustime",
		},
		{
			Name:  "probe_ref",
			Value: ep.ProbeRef.Id(),
		},
		{
			Name:  "probe",
			Value: ep.ProbeRef.Probe,
		},
		{
			Name:  "group",
			Value: ep.ProbeRef.Group,
		},
	}
}

func stringifyTimeseries(tss []*prompb.TimeSeries, name string) string {
	b := strings.Builder{}
	for _, ts := range tss {
		b.WriteString("\n" + name + "   ")
		b.WriteString(stringifyLabels(ts.Labels))
		for _, s := range ts.Samples {
			stamp := time.Unix(s.Timestamp/1000, 0).Format("15:04:05")
			b.WriteString(fmt.Sprintf("    %s  %0.f", stamp, s.Value))
		}
	}
	return b.String()
}

func stringifyLabels(labels []*prompb.Label) string {
	var ref, status, name string
	for _, lbl := range labels {
		if lbl.Name == "probe_ref" {
			ref = lbl.Value
			continue
		}
		if lbl.Name == "status" {
			status = lbl.Value
		}

		if lbl.Name == "__name__" {
			name = lbl.Value
		}
	}
	return fmt.Sprintf("__name__=%s ref=%s status=%s", name, ref, status)
}