#!/bin/bash

cp ./tmp/generated-gateway_config.yml ./tmp/openfaas-cloud/gateway_config.yml
cp ./tmp/generated-github.yml ./tmp/openfaas-cloud/github.yml

cd ./tmp/openfaas-cloud


export ADMIN_PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'|base64 -D)
faas-cli template pull

kubectl port-forward svc/gateway -n openfaas 31111:8080 &
sleep 2

export OPENFAAS_URL=http://127.0.0.1:31111
echo -n $ADMIN_PASSWORD | faas-cli login --username admin --password-stdin

faas-cli deploy

kill %1
