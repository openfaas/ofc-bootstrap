#!/bin/bash

kubectl apply \
    -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml

kubectl create namespace cert-manager
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true

helm install \
    --name cert-manager \
    --namespace cert-manager \
    --version v0.6.0 \
    stable/cert-manager
