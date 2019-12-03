#!/bin/bash
./scripts/reset-kind.sh

export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"

cd $GOPATH/src/github.com/openfaas-incubator/ofc-bootstrap

# Fake the secrets from init.yaml
mkdir -p ~/Downloads
mkdir -p ~/.docker/ 
touch ~/Downloads/secret-access-key
touch ~/Downloads/private-key.pem
touch ~/Downloads/do-access-token
touch ~/.docker/config.json

# Run end to end
./bin/ofc-bootstrap apply --file example.init.yaml
