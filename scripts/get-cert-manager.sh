#!/bin/bash

# cert-manager is ready for CRD objects when this condition is "True"
kubectl rollout status cert/cert-manager-webhook-webhook-tls -n  cert-manager --timeout=10m
kubectl rollout status deploy/cert-manager-webhook -n cert-manager --timeout=10m