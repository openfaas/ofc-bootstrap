#!/bin/bash
./scripts/reset-kind.sh
export KUBECONFIG="$(kind get kubeconfig-path --name="1")"
go run main.go -yaml init.yaml

