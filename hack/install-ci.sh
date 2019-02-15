#!/bin/bash

curl -sLSf https://get.docker.com | sudo sh

curl -OSL https://dl.google.com/go/go1.10.7.linux-amd64.tar.gz
mkdir -p /usr/local/go
tar -xvf go1.10.7.linux-amd64.tar.gz --strip-components=1 -C /usr/local/go/

echo "export GOPATH=\$HOME/go" | tee -a ~/.bash_profile
echo "export PATH=\$GOPATH/bin:\$PATH:/usr/local/go/bin/" | tee -a ~/.bash_profile

curl -sLSf https://cli.openfaas.com | sudo sh

sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

curl -sLSf https://raw.githubusercontent.com/helm/helm/master/scripts/get | sudo bash

go get sigs.k8s.io/kind
