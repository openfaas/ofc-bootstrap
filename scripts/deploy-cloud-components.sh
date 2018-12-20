#!/bin/bash

rm -rf ./tmp/generated-*
cp ./tmp/generated-gateway_config.yml ./tmp/openfaas-cloud/gateway_config.yml
cp ./tmp/generated-github.yml ./tmp/openfaas-cloud/github.yml
cp ./tmp/generated-dashboard_config.yml ./tmp/openfaas-cloud/dashboard/dashboard_config.yml

cd ./tmp/openfaas-cloud


export ADMIN_PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 --decode)
faas-cli template pull 

kubectl port-forward svc/gateway -n openfaas 31111:8080 &
sleep 2


while [ true ];
do
    echo "Checking if OpenFaaS GW is up."
    curl -if 127.0.0.1:31111
    if [ $? == 0 ];
    then
        break
    fi
done


export OPENFAAS_URL=http://127.0.0.1:31111
echo -n $ADMIN_PASSWORD | faas-cli login --username admin --password-stdin

faas-cli deploy

cd ./dashboard
faas-cli template pull https://github.com/openfaas-incubator/node8-express-template
faas-cli deploy

kill %1
