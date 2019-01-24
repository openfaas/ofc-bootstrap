#!/bin/bash
./scripts/reset-kind.sh
export KUBECONFIG="$(kind get kubeconfig-path --name="1")"
go run main.go -yaml init.yaml


# Fake the secrets from init.yaml

mkdir -p ~/Downloads
mkdir -p ~/.docker/ 
touch ~/Downloads/secret-access-key
touch ~/Downloads/ofc-bootstrap-test.2018-12-23.private-key.pem
touch ~/.docker/config.json