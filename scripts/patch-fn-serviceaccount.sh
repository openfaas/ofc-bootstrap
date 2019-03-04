#!/bin/bash

kubectl patch serviceaccount default -p '{"secrets": [{"name": "registry-secret"}]}' -n openfaas-fn