#!/bin/bash

helm delete --purge cert-manager nginxingress openfaas cloud-minio
kubectl delete ns openfaas openfaas-fn
kubectl delete crd sealedsecrets.bitnami.com
kubectl delete deploy/sealed-secrets-controller -n kube-system
kubectl delete deploy/tiller-deploy -n kube-system
kubectl delete sa/tiller -n kube-system
kubectl delete clusterrolebinding/tiller -n kube-system
kubectl delete certificates --all  -n openfaas

rm -rf ./tmp
