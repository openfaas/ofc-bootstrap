#!/bin/bash

curl -sLSf https://get.docker.com | sudo sh

echo "export GOPATH=\$HOME/go" | tee -a ~/.bash_profile
echo "export PATH=\$GOPATH/bin:\$PATH:/usr/local/go/bin/" | tee -a ~/.bash_profile

curl -sLSf https://cli.openfaas.com | sudo sh

curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

curl -sLSf https://raw.githubusercontent.com/helm/helm/master/scripts/get | sudo bash

curl -SLfs https://github.com/kubernetes-sigs/kind/releases/download/v0.5.1/kind-linux-amd64 > kind-linux-amd64
chmod +x kind-linux-amd64 
sudo mv kind-linux-amd64 /usr/local/bin/kind

kind create cluster
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"

# wait, roughly for the cluster to finish starting
kubectl rollout status deploy coredns --watch -n kube-system
