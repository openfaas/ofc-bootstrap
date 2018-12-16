#!/bin/bash

kubectl get deploy/tiller-deploy -n kube-system --output="jsonpath={.status.availableReplicas}"
