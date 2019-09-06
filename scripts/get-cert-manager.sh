#!/bin/bash

# cert-manager is ready for CRD objects when this condition is "True"
INJECTOR_READY=$(kubectl get deploy/cert-manager-cainjector -n  cert-manager -o jsonpath="{.status.conditions[0].status}")
CERT_MANAGER_READY=$(kubectl get deploy/cert-manager -n cert-manager -o jsonpath="{.status.conditions[0].status}")
CERT_MANAGER_WEBHOOK_READY=$(kubectl get deploy/cert-manager-webhook -n cert-manager -o jsonpath="{.status.conditions[0].status}")

if [ "$CERT_MANAGER_READY" = "True" ]
then
    if [ "$INJECTOR_READY" = "True" ]
    then
        if [ "$CERT_MANAGER_WEBHOOK_READY" = "True" ]
        then
            echo -n True
        fi
    fi
fi
