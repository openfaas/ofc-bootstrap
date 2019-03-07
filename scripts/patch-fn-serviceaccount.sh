#!/bin/bash

kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "registry-pull-secret"}]}' -n openfaas-fn