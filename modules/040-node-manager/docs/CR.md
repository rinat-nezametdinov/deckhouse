---
title: "Управление узлами: custom resources"
---

## NodeGroup

Описывает runtime параметры группы узлов

Все опции идут в `.spec`.

* `nodeType` — тип узлов, которые представляет эта группа.
  * Доступны следующие значения:
      * Cloud — узлы для этой группы будут автоматически создаваться (и удаляться) в настроенном облачном провайдере,
      * Static — статический узел, размещенный на железном сервере или виртуальной машине. Узел не управляется
        cloud-controller-manager'ом, даже если включен один из облачных провайдеров.
      * Hybrid – статический узел (созданный вручную или любыми внешними инструментами), размещенный в том же облаке, с
        которым настроена интеграция у одного из облачных провайдеров, на таком узле работает CSI и такой узел
        управляется cloud-controller-manager'ом (объект Node автоматически обогащается информацией о зоне и регионе по
        данным, полученным от облака; при удалении узла из облака, соответствующий ему Node-объект будет
        удален в Kubernetes).
* `disruptions` — параметры обновлений приводящих к возможному простою.
    * `approvalMode` — режим выдачи разрешения на disruptive обновление.
        * `Automatic` —  автоматически выдавать разрешения на disruption при обновлении (значение по умолчанию).
        * `Manual` — отключить автоматическую выдачу разрешений на disruption, когда disruption потребуется – загарится специальный алерт.
    * `automatic` — дополнительные параметры для режима `Automatic`.
        * `drainBeforeApproval` — выгон (draining) подов с ноды, перед выдачей разрешения на disruption.
            * Формат — boolean.
            * По умолчанию, `true`.
* `kubernetesVersion` — желаемая minor-версия Kubernetes.
  * Например, `1.16`.
  * По умолчанию соответствует глобально выбранной для кластера версии (см. документацию по установке) или, если таковая не определена, текущей версии control-plane'а.
* `static` — параметры связанные со статическими узлами.
  * **Внимание!** Допустимо использовать только совместно с `nodeType: Static`
  * `internalNetworkCIDRs` — список подсетей, использующиеся для коммуникации внутри кластера. На основании этого списка
    производится автоматическое определение InternalIP узла и адреса, на котором будет слушать kubelet.
    * Формат — массив строк. Subnet CIDR.
    * Пример:

      ```yaml
      internalNetworkCIDRs:
      - "10.2.2.3/24"
      - "10.1.1.1/24"
      ```
* `cloudInstances` – параметры заказа облачных виртуальных машин.
  * **Внимание!** Допустимо использовать только совместно с `nodeType: Cloud`
  * `classReference` – ссылка на объект InstanceClass. Уникален для каждого `cloud-provider-` модуля.
    * `kind` — тип объекта (например, `OpenStackInstanceClass`). Тип объекта указан в документации соответствующего
      `cloud-provider-` модуля.
    * `name` — имя нужного InstanceClass объекта (например, `finland-medium`).
  * `maxPerZone` — максимальное количество инстансов в зоне. Проставляется как верхняя граница в cluster-autoscaler.
  * `minPerZone` — минимальное количество инстансов в зоне. Проставляется в объект MachineDeployment и в качестве нижней
     границы в cluster-autoscaler.
    * **Внимание!** Не может быть меньше 1.
  * `maxUnavailablePerZone` — сколько инстансов может быть недоступно при RollingUpdate'е.
    * По умолчанию `0`.
  * `maxSurgePerZone` — сколько инстансов создавать одновременно при scale-up.
    * По умолчанию `1`.
  * `zones` — переопределение перечня зон, в которых создаются инстансы.
    * Формат — массив строк.
    * Опциональный параметр.
    * Значение по умолчанию зависит от выбранного облачного провайдера и обычно соответствует всем зонам используемого
      региона.
  * `standby` — количество избыточно выделенных узлов для этой `NodeGroup`:
    * Опциональный параметр.
    * Значение может быть абсолютным (например, `2`) или процентом желаемых узлов (например, `10%`).
    * Абсолютное значение рассчитывается из процента от максимального количества узлов путем округления в меньшую сторону, но минимум — `1`.
* `operatingSystem` — параметры операционной системы.
  * `manageKernel` — автоматическое управление ядром операционной системы.
    * Формат — boolean.
    * По умолчанию, `true`.
* `kubelet` — параметры настройки kubelet'а.
  * `maxPods` — максимальное количество подов на нодах данной `NodeGroup`.
    * По умолчанию `110`.
  * `rootDir` — Путь к каталогу для файлов kubelet'а (volume mounts, ...).
    * По умолчанию `/var/lib/kubelet`.
* `docker` — параметры настройки docker'а.
  * `maxConcurrentDownloads` — максимальное количество потоков одновременного скачивания docker образов.
    * По умолчанию `3`.
  * `manage` — автоматическое управление версией и параметрами docker.
    * По умолчанию `true`.
* `nodeTemplate` — настройки Node объектов в Kubernetes, которые будут добавлены после регистрации ноды.
  * `labels` — аналогично стандартному [полю](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#objectmeta-v1-meta) `metadata.labels`
    * Пример:

      ```yaml
      labels:
        environment: production
        app: warp-drive-ai

  * `annotations` — аналогично стандартному [полю](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#objectmeta-v1-meta) `metadata.annotations`
    * Пример:

      ```yaml
      annotations:
        ai.fleet.com/discombobulate: "true"
      ```

  * `taints` — аналогично полю `.spec.taints` из объекта [Node](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#taint-v1-core). **Внимание!** Доступны только поля `effect`, `key`, `values`.
    * Пример:

      ```yaml
      taints:
      - effect: NoExecute
        key: ship-class
        value: frigate
      ```

* `chaos` — настройки chaos monkey:
  * Опциональный параметр.
  * `mode` — режим работы chaos monkey, возможные значения: `DrainAndDelete` — при срабатывании drain'ит и удаляет ноду, `Disabled` — не трогает данную NodeGroup.
    * По умолчанию `Disabled`.
  * `period` — в какой интервал времени сработает chaos monkey (указывать можно в [golang формате](https://golang.org/pkg/time/#ParseDuration));
    * По умолчанию `6h`.

## NodeUser

Описывает linux-пользователей, которые будут созданы на всех узлах.

Все опции расположены в поле `.spec`.

* `uid` — user id пользователя на узлах.
  * Формат — число > 1000.
  * Обязательный параметр.
  * Неизменяемый в течение жизни ресурса параметр.
* `sshPublicKey` — публичный ssh ключ пользователя.
  * Формат соответствует содержимому файла публичного файла ключа, например — `ssh-rsa AAAAB3NzaC1yc2EAAA...`
  * Обязательный параметр.
* `passwordHash` — хеш пароля пользователя.
  * Формат соответствует хешам паролей, содержащихся в `/etc/shadow`, например — `$2a$10$GAwx2h0D1...`, можно получить при помощи команды `openssl passwd -6`.
  * Обязательный параметр.
* `isSudoer` — определяет, будет ли пользователю разрешено sudo.
  * Формат — boolean (true/false).
  * Необязательный параметр.
  * По умолчанию — false.
* `extraGroups` — список дополнительных групп, в которые должен быть включен пользователь.
  * Формат — список строк.
  * Необязательный параметр.