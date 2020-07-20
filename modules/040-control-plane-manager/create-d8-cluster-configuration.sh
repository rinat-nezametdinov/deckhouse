#!/bin/bash

set -ex

provider=""
clusterType="Cloud"
p=$(deckhouse-controller module list | grep cloud-provider | cut -d- -f3 | sed 's/ *$//g')

case "$p" in
  'aws')
    provider="AWS"
    ;;
  'gcp')
    provider="GCP"
    ;;
  'openstack')
    provider="OpenStack"
    ;;
  'yandex')
    provider="Yandex"
    ;;
  'vsphere')
    provider="Vsphere"
    ;;
  * )
    clusterType="Static"
    ;;
esac

#check if kops
if kubectl -n kube-system get pods -l k8s-app=kube-controller-manager -o yaml | grep 'cloud-provider' > /dev/null 2> /dev/null ; then
  echo "Error: script must not be executed on kops clusters"
  exit 1
fi

#check if single master
if [[ "x$(kubectl get nodes -l node-role.kubernetes.io/master="" -o name | wc -l)" != "x1" ]] ; then
  echo "Error: script must not be executed on multi master clusters"
  exit 1
fi

#check if d8-pki exists
if ! kubectl -n kube-system get secret d8-pki > /dev/null 2> /dev/null ; then
  echo "Error: d8-pki secret in namespace kube-system doesn't exists"
  exit 1
fi

masterInternalIP=$(kubectl get nodes -l node-role.kubernetes.io/master='' -o json | jq '.items[0].status.addresses[] | select(.type == "InternalIP") | .address' -r)

#check etcd on correct ip
if ! nc -z "$masterInternalIP" 2379; then
  echo "Error: etcd must listen on master internal ip"
  exit 1
fi

apiserverAdvertiseAddress=$(kubectl -n default get ep kubernetes -o json | jq '.subsets[0].addresses[0].ip' -r)
#check apiserver on correct ip and controlPlaneManager configuration
if [[ "x$apiserverAdvertiseAddress" != "x$masterInternalIP" ]]; then
  if [[ "x$(deckhouse-controller module values control-plane-manager | yq r - controlPlaneManager.apiserver.bindToWildcard)" != "xtrue" ]]; then
    echo "Error: apiserver advertise address differs from master internal ip, so you must set controlPlaneManager.apiserver.bindToWildcard to true"
    exit 1
  fi

  if ! deckhouse-controller module values control-plane-manager | yq r - controlPlaneManager.apiserver.certSANs -j | jq -e ". | index(\"$apiserverAdvertiseAddress\")"; then
    echo "Error: apiserver advertise address must be added to controlPlaneManager.apiserver.certSANs"
    exit 1
  fi
fi

#check if kubeadm config exists
if ! kubectl -n kube-system get cm kubeadm-config > /dev/null 2> /dev/null ; then
  echo "Error: kubeadm-config cm in namespace kube-system doesn't exists"
  exit 1
fi

cluster_configuration=''
if [[ $clusterType == "Static" ]]; then
  cluster_configuration=$(cat <<END
apiVersion: deckhouse.io/v1alpha1
kind: ClusterConfiguration
clusterType: $clusterType
podSubnetCIDR: $(deckhouse-controller global values | yq r - global.discovery.podSubnet)
podSubnetNodeCIDRPrefix: $(kubectl -n kube-system get cm kubeadm-config -o yaml | yq r - data.ClusterConfiguration | yq r -j - controllerManager.extraArgs | jq '."node-cidr-mask-size" //= "24" | ."node-cidr-mask-size"')
serviceSubnetCIDR: $(deckhouse-controller global values | yq r - global.discovery.serviceSubnet)
kubernetesVersion: "$(deckhouse-controller global values | yq r - global.discovery.kubernetesVersion | cut -c 1-4)"
clusterDomain: $(deckhouse-controller global values | yq r - global.discovery.clusterDomain)
END
)
elif [[ $clusterType == "Cloud" ]]; then
  cluster_configuration=$(cat <<END
apiVersion: deckhouse.io/v1alpha1
kind: ClusterConfiguration
clusterType: $clusterType
cloud:
  provider: $provider
  prefix: $(deckhouse-controller module values node-manager | yq r - nodeManager.instancePrefix)
podSubnetCIDR: $(deckhouse-controller global values | yq r - global.discovery.podSubnet)
podSubnetNodeCIDRPrefix: $(kubectl -n kube-system get cm kubeadm-config -o yaml | yq r - data.ClusterConfiguration | yq r -j - controllerManager.extraArgs | jq '."node-cidr-mask-size" //= "24" | ."node-cidr-mask-size"')
serviceSubnetCIDR: $(deckhouse-controller global values | yq r - global.discovery.serviceSubnet)
kubernetesVersion: "$(deckhouse-controller global values | yq r - global.discovery.kubernetesVersion | cut -c 1-4)"
clusterDomain: $(deckhouse-controller global values | yq r - global.discovery.clusterDomain)
END
)
fi

cat <<END
##### check cluster_configuration.yaml#######
$cluster_configuration
#############################################
END

if [[ $1 == "dry-run" ]]; then
  exit 0
fi

kubectl create -f - <<END
apiVersion: v1
data:
  cluster-configuration.yaml: $(echo "$cluster_configuration" | base64 -w 0)
kind: Secret
metadata:
  labels:
    heritage: deckhouse
  name: d8-cluster-configuration
  namespace: kube-system
type: Opaque
END