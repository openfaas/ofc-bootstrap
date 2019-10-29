#!/bin/bash

helm repo add openfaas https://openfaas.github.io/faas-netes

helm repo update && \
helm upgrade openfaas --install openfaas/openfaas \
    --namespace openfaas  \
    --set basic_auth=true \
    --set functionNamespace=openfaas-fn \
    --set ingress.enabled=true \
    --set gateway.scaleFromZero=true \
    --set gateway.readTimeout=15m \
    --set gateway.writeTimeout=15m \
    --set gateway.upstreamTimeout=14m55s \
    --set queueWorker.ackWait=15m \
    --set faasnetes.readTimeout=5m \
    --set faasnetes.writeTimeout=5m \
    --set gateway.replicas=2 \
    --set queueWorker.replicas=2 \
    --set faasIdler.dryRun=$FAAS_IDLER_DRY_RUN \
    --set faasnetes.httpProbe=true \
    --set faasnetes.imagePullPolicy=IfNotPresent
