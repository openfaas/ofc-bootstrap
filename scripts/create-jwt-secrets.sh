#!/bin/bash

# Remove old keys
rm -rf ./key ./key.pub

# Private key
openssl ecparam -genkey -name prime256v1 -noout -out key

# Public key
openssl ec -in key -pubout -out key.pub

kubectl -n openfaas create secret generic jwt-private-key --from-file=./key
kubectl -n openfaas create secret generic jwt-public-key --from-file=./key.pub