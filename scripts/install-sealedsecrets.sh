#!/bin/bash

release=$(curl -sI https://github.com/bitnami-labs/sealed-secrets/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
#release=$(curl --silent "https://api.github.com/repos/bitnami-labs/sealed-secrets/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')

echo "SealedSecrets release: $release"

helm del --purge ofc-sealedsecrets
kubectl delete customresourcedefinition sealedsecrets.bitnami.com

helm install --namespace kube-system --name ofc-sealedsecrets stable/sealed-secrets