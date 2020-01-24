#!/bin/bash

set -e

helm upgrade --install \
 --namespace kube-system ofc-sealedsecrets stable/sealed-secrets \
  --wait

kubectl rollout status deploy/ofc-sealedsecrets-sealed-secrets -n kube-system --timeout 5m
