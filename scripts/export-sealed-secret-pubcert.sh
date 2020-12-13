#!/bin/bash
kubeseal --version
kubeseal --fetch-cert --controller-name=sealed-secrets > tmp/pub-cert.pem  \
  && cat tmp/pub-cert.pem
