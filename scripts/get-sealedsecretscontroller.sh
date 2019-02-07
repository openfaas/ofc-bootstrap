#!/bin/bash

kubectl get deploy/ofc-sealedsecrets-sealed-secrets -n kube-system --output="jsonpath={.status.availableReplicas}"
