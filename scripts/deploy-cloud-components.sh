#!/bin/bash

cp ./tmp/generated-gateway_config.yml ./tmp/openfaas-cloud/gateway_config.yml
cp ./tmp/generated-github.yml ./tmp/openfaas-cloud/github.yml

cd ./tmp/openfaas-cloud && faas-cli template pull && faas-cli deploy