#!/bin/bash

set -e

export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"

cd $GOPATH/src/github.com/openfaas/ofc-bootstrap

# Fake the secrets from init.yaml
mkdir -p ~/Downloads
touch ~/Downloads/secret-access-key
touch ~/Downloads/private-key.pem
touch ~/Downloads/do-access-token

# Run end to end

./bin/ofc-bootstrap registry-login --username fake --password also-fake
./bin/ofc-bootstrap apply --file example.init.yaml
