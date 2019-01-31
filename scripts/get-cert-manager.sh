#!/bin/bash

# cert-manager is ready for CRD objects when this condition is "True"
CERT_READY=$(kubectl get cert/cert-manager-webhook-webhook-tls -n  cert-manager -o jsonpath="{.status.conditions[0].status}")
WEBHOOK_READY=$(kubectl get deploy/cert-manager-webhook -n cert-manager -o jsonpath="{.status.conditions[0].status}")

if [ "$CERT_READY" = "True" ]
then
    if [ "$WEBHOOK_READY" = "True" ]
    then
        echo -n True
    fi
fi
