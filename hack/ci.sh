#!/bin/bash
./scripts/reset-kind.sh
export KUBECONFIG="$(kind get kubeconfig-path --name="1")"

#add kubectl in ci stage
sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

# Build the code
go build

# Fake the secrets from init.yaml
mkdir -p ~/Downloads
mkdir -p ~/.docker/ 
touch ~/Downloads/secret-access-key
touch ~/Downloads/private-key.pem
touch ~/Downloads/do-access-token
touch ~/.docker/config.json

# Run end to end
./ofc-bootstrap -yaml example.init.yaml
