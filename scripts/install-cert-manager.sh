#!/bin/bash

helm install \
    --name cert-manager \
    --namespace kube-system \
    --version v0.4.0 \
    stable/cert-manager