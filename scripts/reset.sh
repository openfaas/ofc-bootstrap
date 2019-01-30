#!/bin/bash

helm delete --purge cert-manager nginxingress openfaas cloud-minio

kubectl delete certificates --all -n openfaas
kubectl delete clusterissuer letsencrypt-prod letsencrypt-staging

kubectl delete ns openfaas openfaas-fn cert-manager

kubectl delete crd sealedsecrets.bitnami.com
kubectl delete \
    -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml
kubectl delete deploy/sealed-secrets-controller -n kube-system
kubectl delete deploy/tiller-deploy -n kube-system
kubectl delete sa/tiller -n kube-system
kubectl delete clusterrolebinding/tiller -n kube-system

kubectl delete secret/clouddns-service-account -n kube-system

rm -rf ./tmp
