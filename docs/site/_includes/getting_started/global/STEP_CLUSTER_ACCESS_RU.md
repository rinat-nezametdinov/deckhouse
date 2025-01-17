<script type="text/javascript" src='{{ assets["getting-started.js"].digest_path }}'></script>
<script type="text/javascript" src='{{ assets["getting-started-access.js"].digest_path }}'></script>

# Доступ к кластеру через Kubernetes API
Deckhouse только что завершил процесс установки вашего кластера. Теперь вы можете подключиться к мастер-узлу, используя ssh.
Для этого необходимо получить IP-адрес мастера либо из логов dhctl, либо из web интерфейса/cli утилиты облачного провайдера.
{% snippetcut %}
```shell
ssh {% if page.platform_code == "azure" %}azureuser{% elsif page.platform_code == "gcp" %}user{% else %}ubuntu{% endif %}@<MASTER_IP>
```
{% endsnippetcut %}
Вы можете запускать kubectl на мастере от пользователя root. Это не безопасный способ, и мы рекомендуем настроить [внешний доступ](/ru/documentation/v1/modules/150-user-authn/usage.html#внешний-доступ-к-kubernetes-api) к Kubernetes API позже.
{% snippetcut %}
```shell
sudo -i
kubectl get nodes
```
{% endsnippetcut %}

# Доступ к кластеру через NGINX Ingress
[IngressNginxController](/en/documentation/v1/modules/402-ingress-nginx/cr.html#ingressnginxcontroller) был создан во время процесса установки кластера.
Теперь осталось настроить доступ к веб-интерфейсам компонентов, которые уже установлены в кластере, таким как Grafana, Prometheus, Dashboard и так далее.
LoadBalancer уже создан и вам остаётся только направить DNS-домен на него.
В первую очередь необходимо подключиться к мастер-узлу, как это описано [выше](#доступ-к-кластеру-через-kubernetes-api).

{% if page.platform_type == 'cloud' %}
Получите IP адрес балансировщика. Для этого в кластере от пользователя root выполните команду:
{% if page.platform_code == 'aws' %}
{% snippetcut %}
{% raw %}
```shell
BALANCER_HOSTNAME=$(kubectl -n d8-ingress-nginx get svc nginx-load-balancer -o json | jq -r '.status.loadBalancer.ingress[0].hostname')
echo "$BALANCER_HOSTNAME"
```
{% endraw %}
{% endsnippetcut %}
{% else %}
{% snippetcut %}
{% raw %}
```shell
BALANCER_IP=$(kubectl -n d8-ingress-nginx get svc nginx-load-balancer -o json | jq -r '.status.loadBalancer.ingress[0].ip')
echo "$BALANCER_IP"
```
{% endraw %}
{% endsnippetcut %}
{% endif %}
{% endif %}

Настройте домен для сервисов Deckhouse, который вы указали на шаге «[Установка кластера](./step3.html)», одним из следующих способов:
<div markdown="1">
<ul><li><p>Если у вас есть возможность добавить DNS-запись используя DNS-сервер, то мы рекомендуем добавить
{%- if page.platform_code == 'aws' %} wildcard CNAME-запись для <code>*.example.com</code> со значением адреса балансировщика (<code>BALANCER_HOSTNAME</code>)
{%- else %} wildcard A-запись для <code>*.example.com</code> со значением IP-адреса балансировщика (<code>BALANCER_IP</code>)
{%- endif -%}
  , который вы получили выше.</p></li>
<li><p>Если вы не имеете под управлением DNS-сервер, добавьте статические записи в файл <code>/etc/hosts</code> для Linux (<code>%SystemRoot%\system32\drivers\etc\hosts</code> для Windows).</p>
{% if page.platform_code == 'aws' %}
  <p>Определить IP-адрес балансировщика можно при помощи следующей команды (также выполняемой в кластере):</p>

<div markdown="1">
{% snippetcut %}
```bash
BALANCER_IP=$(dig "$BALANCER_HOSTNAME" +short | head -1); echo "$BALANCER_IP"
```
{% endsnippetcut %}
</div>
{% endif %}

  <p>Для добавления записей в файл <code>/etc/hosts</code> локально, выполните например следующие шаги:</p>

<ul><li><p>Экспортируйте переменную <code>BALANCER_IP</code>, указав полученный IP-адрес балансировщика:</p>
{% snippetcut %}
```bash
export BALANCER_IP="<PUT_BALANCER_IP_HERE>"
```
{% endsnippetcut %}
</li>
  <li><p>Добавьте DNS-записи для веб-интерфейсов Deckhouse:</p>
{% snippetcut selector="example-hosts" %}
```bash
sudo -E bash -c "cat <<EOF >> /etc/hosts
$BALANCER_IP dashboard.example.com
$BALANCER_IP deckhouse.example.com
$BALANCER_IP kubeconfig.example.com
$BALANCER_IP grafana.example.com
$BALANCER_IP dex.example.com
EOF
"
```
{% endsnippetcut %}
</li>
</ul></li>
</ul>
</div>
