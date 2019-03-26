#!/bin/bash

kubectl rollout status deploy/ofc-sealedsecrets-sealed-secrets -n kube-system
