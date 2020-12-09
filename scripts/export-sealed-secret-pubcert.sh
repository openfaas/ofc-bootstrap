#!/bin/bash

if [ ! -f /tmp/kubeseal ];
then
    GOOS=$(echo $(uname -s) | tr '[:upper:]' '[:lower:]')
    GOARCH=$(echo $(uname -m) | tr '[:upper:]' '[:lower:]')

    if [ "$GOARCH" == "x86_64" ]; then
       GOARCH="amd64"
    fi

    faas-cli cloud seal --download --download-to=/tmp/
    chmod +x /tmp/kubeseal
fi

/tmp/kubeseal --fetch-cert --controller-name=ofc-sealedsecrets-sealed-secrets > tmp/pub-cert.pem && \
  cat tmp/pub-cert.pem
