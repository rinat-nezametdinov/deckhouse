---
title: "Управление control plane: FAQ"
---

## Как добавить мастер?

Поставить на узле кластер label `node-role.kubernetes.io/master: ""`, все остальное произойдет полностью автоматически.

## Как удалить мастер?

* Если удаление не нарушает кворум в etcd (в корректно функционирующем кластере это все ситуации, кроме перехода 2 -> 1):
    1. Если удаление мастера согласовано, то необходимо удалить виртуальную машину обычным способом.
       Если удалить мастер прямо сейчас нельзя, например, на нём настроены бэкапы или деплой выполняется через него, то необходимо остановить используемый на узле Container Runtime
       
       В случае использования Docker:
       ```
       systemctl stop docker
       systemctl disable docker
       ```
       В случае использования Containerd:
       ```
       systemctl stop containerd
       systemctl disable containerd
       ```
    2. Удалить объект Node из Kubernetes
    3. [Дождаться](#как-посмотреть-список-memberов-в-etcd), пока etcd member будет автоматически удален.

* Если удаление нарушает кворум (переход 2 -> 1):
    1. Остановить kubelet на узле (не останавливая контейнер с etcd):
       ```
       systemctl stop kubelet
       systemctl stop bashible.timer
       systemctl stop bashible
       systemctl disable kubelet
       systemctl disable bashible.timer
       systemctl disable bashible
       ```
    2. Удалить объект Node из Kubernetes
    3. [Дождаться](#как-посмотреть-список-memberов-в-etcd), пока etcd member будет автоматически удален.
    4. Если удаление мастера согласовано, то необходимо удалить виртуальную машину обычным способом.
       Если удалить мастер прямо сейчас нельзя, например, на нём настроены бэкапы или деплой выполняется через него, то необходимо остановить используемый на ноде Container Runtime
       
       В случае использования Docker:
       ```
       systemctl stop docker
       systemctl disable docker
       ```
       В случае использования Containerd:
       ```
       systemctl stop containerd
       systemctl disable containerd
       ```

## Как убрать мастер, сохранив узел?

1. Снять лейблы `node.deckhouse.io/group: master` и `node-role.kubernetes.io/master: ""`, затем дождаться, пока etcd member будет удален автоматически.
2. Зайти на узел и выполнить следующие действия:
  ```shell
  rm -f /etc/kubernetes/manifests/{etcd,kube-apiserver,kube-scheduler,kube-controller-manager}.yaml
  rm -f /etc/kubernetes/{scheduler,controller-manager}.conf
  rm -f /etc/kubernetes/authorization-webhook-config.yaml
  rm -f /etc/kubernetes/admin.conf /root/.kube/config
  rm -rf /etc/kubernetes/deckhouse
  rm -rf /etc/kubernetes/pki/{ca.key,apiserver*,etcd/,front-proxy*,sa.*}
  rm -rf /var/lib/etcd
  ```

## Как посмотреть список member'ов в etcd?

1. Зайти в pod с etcd.
  ```shell
  kubectl -n kube-system exec -ti $(kubectl -n kube-system get pod -l component=etcd,tier=control-plane -o name | head -n1) sh
  ```
2. Выполнить команду.
  ```shell
  ETCDCTL_API=3 etcdctl --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/ca.crt --key /etc/kubernetes/pki/etcd/ca.key --endpoints https://127.0.0.1:2379/ member list
  ```

## Что делать, если что-то пошло не так?

В процессе работы control-plane-manager оставляет резервные копии в `/etc/kubernetes/deckhouse/backup`, они могут помочь.

## Что делать, если кластер etcd развалился?

1. Остановить (удалить `/etc/kubernetes/manifests/etcd.yaml`) etcd на всех узлах, кроме одной. С неё мы начнём восстановление multi-master'а.
2. На оставшемся узле указать следующий параметр командной строки: `--force-new-cluster`.
3. После успешного подъёма кластера, удалить параметр `--force-new-cluster`.

**Внимание!** Операция деструктивна, полностью уничтожает консенсус и запускает etcd кластер с состояния, которое сохранилось на узле. Любые pending записи пропадут.

## Как включить аудит событий?

Если вам требуется вести учёт операций в кластере или отдебажить неожиданное поведение - для всего этого в Kubernetes предусмотрен [Auditing](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-cluster/), настраиваемый через указание соответствующих [Audit Policy](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/#audit-policy).

На данный момент используются фиксированные параметры ротации логов:
```bash
--audit-log-maxage=7
--audit-log-maxbackup=10
--audit-log-maxsize=100
```
Предполагается наличие на master-узлах `скрейпера логов` *(filebeat, promtail)*, который будет следить за директорией с логами:
```bash
/var/log/kube-audit/audit.log
```
В зависимости от настроек `Policy` и количества запросов к **apiserver** - логов может быть очень много, соответственно глубина хранения может быть менее 30 минут. Максимальное занимаемое место на диске ограничено `1000 Мб`. Логи старше `7 дней` так же будут удалены.

### Предостережение
> ⚠️ Текущая реализация функционала не является безопасной, с точки зрения, возможности временно сломать **control-plane**.
>
> Если в секрете с конфигурационным файлом окажутся неподдерживаемые опции или опечатка, то **apiserver** не сможет запуститься.

В случае возникновения проблем с запуском **apiserver** потребуется вручную отключить параметры `--audit-log-*` в манифесте `/etc/kubernetes/manifests/kube-apiserver.yaml` и перезапустить **apiserver** следующей командой:
```bash
docker stop $(docker ps | grep kube-apiserver- | awk '{print $1}')
```
После перезапуска у вас будет достаточно времени исправить `Secret` или [удалить его](#полезные-команды).

### Включение и настройка
За включение отвечает параметр в `ConfigMap` `d8-system/deckhouse`
```yaml
  controlPlaneManager: |
    apiserver:
      auditPolicyEnabled: true
```
Конфигурация параметров осуществляется через `Secret` `kube-system/audit-policy`, внутрь которого потребуется положить `yaml` файл, закодированный `base64`:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: audit-policy
  namespace: kube-system
data:
  audit-policy.yaml: <base64>
```
### Пример
Минимальный рабочий пример `audit-policy.yaml` выглядит так:
```yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
  omitStages:
  - RequestReceived
```
Подробную информацию по настройке содержимого `audit-policy.yaml` можно получить по следующим ссылкам:
- [В официальной документации Kubernetes](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/#audit-policy)
- [В нашей статье на Habr](https://habr.com/ru/company/flant/blog/468679/)
- [Опираясь на код скрипта-генератора, используемого в GCE](https://github.com/kubernetes/kubernetes/blob/0ef45b4fcf7697ea94b96d1a2fe1d9bffb692f3a/cluster/gce/gci/configure-helper.sh#L722-L862)

### Полезные команды
Создать `Secret` из файла можно командой:
```bash
kubectl -n kube-system create secret generic audit-policy --from-file=./audit-policy.yaml
```
Удалить `Secret` из кластера:
```bash
kubectl -n kube-system delete secret audit-policy
```

## Как ускорить перезапуск подов при потере связи с нодой?

По умолчанию, если нода за 40 секунд не сообщает своё состояние — она помечается как недоступная. И еще через 5 минут её поды начнут перезапускаться на других нодах. Итоговое время недоступности приложений ~ 6 минут.

В специфических случаях, когда приложение не может быть запущено в нескольких экземплярах  — есть способ сократить период их недоступности:  

1. Уменьшить время перехода `Node` в состояние `Unreachable`, при потере с ней связи, настройкой параметра `nodeMonitorGracePeriodSeconds`.
1. Установить меньший таймаут удаления подов с недоступной ноды, в параметре `failedNodePodEvictionTimeoutSeconds`.

### Пример
```yaml
  controlPlaneManager: |
    nodeMonitorGracePeriodSeconds: 10
    failedNodePodEvictionTimeoutSeconds: 50
```
В этом случае при потере связи с нодой, приложения будут перезапущены через ~ 1 минуту.

### Предостережение
Оба описанных параметра оказывают непосредственное влияние на потребляемые `control-plane`'ом  ресурсы процессора и памяти. Уменьшая таймауты, мы заставляем системные компоненты чаще производить отправку статусов и сверки состояний ресурсов. 

В процессе подбора подходящих вам значений обращайте внимание на графики потребления ресурсов управляющих узлов. Будте готовы к тому, что, чем меньшие значения параметров вы выбираете, тем больше ресурсов может потребоваться выделить на эти ноды.      
