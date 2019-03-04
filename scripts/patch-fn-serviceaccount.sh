#!/bin/bash

kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "registry-secret"}]}' -n openfaas-fn