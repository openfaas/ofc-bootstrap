#!/bin/bash

kubectl apply \
    -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml

helm install \
    --name cert-manager \
    --namespace kube-system \
    --version v0.6.0 \
    stable/cert-manager
