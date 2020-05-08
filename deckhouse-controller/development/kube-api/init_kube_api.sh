#!/usr/bin/env bash

apt-get update -y
apt-get -y install iproute2 tcpdump nmap curl apache2-utils
apt-get -y install git jq python3
apt-get -y install gnupg2 libssh2-1-dev libssl-dev cmake

curl -L https://raw.githubusercontent.com/micha/resty/master/resty > /usr/bin/resty
chmod +x /usr/bin/resty
echo "export TERM=xterm" >> /etc/bash.bashrc

#curl -L https://storage.googleapis.com/kubernetes-release/release/v1.7.7/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
#chmod +x /usr/local/bin/kubectl
#gpg2 --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
#curl -sSL https://get.rvm.io | bash -s stable --ruby
#source /etc/profile.d/rvm.sh
#rvm install ruby
#gem install dapp


cat <<DIMG
# Образ для экспериментов с Kubernetes API
dimg_group do
  dimg 'kube-api-tester' do
    docker.from 'ubuntu:16.04'
    shell.before_install do
      run 'apt-get update -y'
      run 'apt-get -y install iproute2 tcpdump nmap curl apache2-utils'
      run 'apt-get -y install git jq python3'
      run 'apt-get -y install gnupg2 libssh2-1-dev libssl-dev cmake'

      run 'curl -L https://raw.githubusercontent.com/micha/resty/master/resty > /usr/bin/resty'
      run 'chmod +x /usr/bin/resty'
      run 'echo "export TERM=xterm" >> /etc/bash.bashrc'

      run 'curl -L https://storage.googleapis.com/kubernetes-release/release/v1.7.7/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl'
      run 'chmod +x /usr/local/bin/kubectl'

      run 'gpg2 --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3'
      run 'curl -sSL https://get.rvm.io | bash -s stable --ruby'
      run 'source /etc/profile.d/rvm.sh'
      run 'rvm install ruby'
      run 'gem install dapp'
    end

    shell.install do
      run 'curl -LOs https://storage.googleapis.com/golang/go1.8.4.linux-amd64.tar.gz'
      run 'tar -C /usr/local -xzf go1.8.4.linux-amd64.tar.gz'
      run 'echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile'
      run 'echo "export GOPATH=/app" >> /etc/profile'
    end

    shell.setup do
      run 'cp /kube_api/kube_api.sh /usr/bin/kube_api'
      run 'chmod +x /usr/bin/kube_api'
    end

    git do
      add '/.kube_api' do
        to '/kube_api'
      end
      add '/' do
        to '/go/src/github.com/deckhouse/deckhouse'
      end
    end
  end
end

DIMG