#!/bin/bash

kubectl rollout status deployment/tiller-deploy -n kube-system
