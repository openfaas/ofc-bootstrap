#!/bin/bash

kubectl apply \
    -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml

helm install \
    --name cert-manager \
    --namespace cert-manager \
    --version v0.6.6 \
    --wait \
    stable/cert-manager
