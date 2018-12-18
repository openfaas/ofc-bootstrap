#!/bin/bash

# Reset script

export KUBECONFIG="$(kind get kubeconfig-path --name="1")"

kind delete cluster --name 1
kind create cluster --name 1
rm ./tmp/*.yml
