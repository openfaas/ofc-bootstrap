#!/bin/bash

kubectl get deploy/sealed-secrets-controller -n kube-system --output="jsonpath={.status.availableReplicas}"
