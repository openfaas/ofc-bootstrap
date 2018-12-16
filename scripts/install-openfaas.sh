#!/bin/bash


helm repo update \
 && helm upgrade openfaas --install openfaas/openfaas \
    --namespace openfaas  \
    --set basic_auth=true \
    --set functionNamespace=openfaas-fn \
    --set ingress.enabled=true \
    --set gateway.scaleFromZero=true \
    --set gateway.readTimeout=300s \
    --set gateway.writeTimeout=300s \
    --set gateway.upstreamTimeout=295s \
    --set faasnetesd.readTimeout=300s \
    --set faasnetesd.writeTimeout=300s \
    --set gateway.replicas=2 \
    --set queueWorker.replicas=2
