#!/bin/bash

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

release=$(curl --silent "https://api.github.com/repos/bitnami-labs/sealed-secrets/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')

curl -sLSf https://github.com/bitnami/sealed-secrets/releases/download/$release/kubeseal-$GOOS-$GOARCH > kubeseal \
chmod +x kubeseal && \
./kubeseal --fetch-cert > pub-cert.pem && \
cat pub-cert.pem
