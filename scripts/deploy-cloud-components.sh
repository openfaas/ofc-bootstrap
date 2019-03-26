#!/bin/bash

cp ./tmp/generated-gateway_config.yml ./tmp/openfaas-cloud/gateway_config.yml
cp ./tmp/generated-github.yml ./tmp/openfaas-cloud/github.yml
cp ./tmp/generated-slack.yml ./tmp/openfaas-cloud/slack.yml
cp ./tmp/generated-dashboard_config.yml ./tmp/openfaas-cloud/dashboard/dashboard_config.yml

kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-builder-dep.yml
kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-builder-svc.yml

kubectl apply -f ./tmp/openfaas-cloud/yaml/core/import-secrets-role.yml

if [ "$ENABLE_OAUTH" = "true" ] ; then
    cp ./tmp/generated-of-auth-dep.yml ./tmp/openfaas-cloud/yaml/core/of-auth-dep.yml
    kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-auth-dep.yml
    kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-auth-svc.yml
    kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-router-dep.yml
else
    #  Disable auth service by pointing the router at the echo function:
    sed s/auth.openfaas/echo.openfaas-fn/g ./tmp/openfaas-cloud/yaml/core/of-router-dep.yml | kubectl apply -f -
fi
kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-router-svc.yml

kubectl apply -f ./tmp/openfaas-cloud/yaml/core/of-auth-svc.yml


cd ./tmp/openfaas-cloud

echo "Creating payload-secret in openfaas-fn"

export PAYLOAD_SECRET=$(kubectl get secret -n openfaas payload-secret -o jsonpath='{.data.payload-secret}'| base64 --decode)

kubectl create secret generic payload-secret -n openfaas-fn --from-literal payload-secret="$PAYLOAD_SECRET"

export ADMIN_PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 --decode)

faas-cli template pull 

kubectl port-forward svc/gateway -n openfaas 31111:8080 &
sleep 2

echo "Checking if OpenFaaS GW is up."

kubectl rollout status deploy/gateway -n openfaas


export OPENFAAS_URL=http://127.0.0.1:31111
echo -n $ADMIN_PASSWORD | faas-cli login --username admin --password-stdin

faas-cli deploy

if [ "$GITLAB" = "true" ] ; then
    cp ../generated-gitlab.yml ./gitlab.yml
    echo "Deploying gitlab functions..."
    faas deploy -f ./gitlab.yml
fi

cd ./dashboard
faas-cli template store pull node10-express
faas-cli deploy

kill %1
