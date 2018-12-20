#!/bin/bash

helm delete --purge cert-manager nginxingress openfaas
kubectl delete ns openfaas openfaas-fn
kubectl delete crd -n kube-system --all
kubectl delete deploy/sealed-secrets-controller -n kube-system
kubectl delete deploy/tiller-deploy -n kube-system
kubectl delete sa/tiller -n kube-system
kubectl delete clusterrolebinding/tiller -n kube-system

