#!/bin/bash

if [ ! -f kubeseal ];
then
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)

    release=$(curl -sI https://github.com/bitnami-labs/sealed-secrets/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')

#    release=$(curl --silent "https://api.github.com/repos/bitnami-labs/sealed-secrets/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
    echo "SealedSecrets release: $release"

    curl -sLSf https://github.com/bitnami/sealed-secrets/releases/download/$release/kubeseal-$GOOS-$GOARCH > kubeseal && \
    chmod +x kubeseal
fi

./kubeseal --fetch-cert > tmp/pub-cert.pem && \
  cat tmp/pub-cert.pem
