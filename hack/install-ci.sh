#!/bin/bash

curl -sLSf https://get.docker.com | sudo sh

curl -OSL https://dl.google.com/go/go1.10.7.linux-amd64.tar.gz
mkdir -p /usr/local/go
tar -xvf go1.10.7.linux-amd64.tar.gz --strip-components=1 -C /usr/local/go/

echo "export GOPATH=\$HOME/go" | tee -a ~/.bash_profile
echo "export PATH=\$GOPATH/bin:\$PATH:/usr/local/go/bin/" | tee -a ~/.bash_profile

curl -sLSf https://cli.openfaas.com | sudo sh

curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

curl -sLSf https://raw.githubusercontent.com/helm/helm/master/scripts/get | sudo bash

export GO111MODULE="on"

go get sigs.k8s.io/kind
