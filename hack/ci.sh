#!/bin/bash
./scripts/reset-kind.sh

export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"

export GO111MODULE=OFF

echo Folder: $(pwd)

cd $GOPATH/src/github.com/openfaas-incubator/ofc-bootstrap

# Build the code
make static

# Fake the secrets from init.yaml
mkdir -p ~/Downloads
mkdir -p ~/.docker/ 
touch ~/Downloads/secret-access-key
touch ~/Downloads/private-key.pem
touch ~/Downloads/do-access-token
touch ~/.docker/config.json

# Run end to end
./ofc-bootstrap -yaml example.init.yaml
