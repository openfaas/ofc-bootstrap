#!/bin/bash
./scripts/reset-kind.sh
export KUBECONFIG="$(kind get kubeconfig-path --name="1")"

# Build the code
go build

# Fake the secrets from init.yaml
mkdir -p ~/Downloads
mkdir -p ~/.docker/ 
touch ~/Downloads/secret-access-key
touch ~/Downloads/private-key.pem
touch ~/.docker/config.json

# Run end to end
./ofc-bootstrap -yaml example.init.yaml
