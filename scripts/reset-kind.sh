#!/bin/bash

kind delete cluster
kind create cluster

export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"

# wait, roughly for the cluster to finish starting
kubectl rollout status deploy coredns --watch -n kube-system

rm ./tmp/*.yml
rm -rf ./tmp/openfaas-cloud
