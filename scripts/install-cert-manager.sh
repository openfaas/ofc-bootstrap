#!/bin/bash

helm install \
    --name cert-manager \
    --namespace kube-system \
    stable/cert-manager
