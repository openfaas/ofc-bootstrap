#!/bin/bash

export USER=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-user}'| base64 --decode)
export PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 --decode)

kubectl create secret generic basic-auth-user \
 --from-literal=basic-auth-user=$USER --namespace openfaas-fn \
 --dry-run -o yaml | kubectl apply -f -

kubectl create secret generic basic-auth-password \
 --from-literal=basic-auth-password=$PASSWORD --namespace openfaas-fn \
 --dry-run -o yaml | kubectl apply -f -
