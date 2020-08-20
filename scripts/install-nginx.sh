#!/bin/bash

set -e

echo helm upgrade --install nginxingress \
  ingress-nginx/ingress-nginx --set rbac.create=true$ADDITIONAL_SET
helm upgrade --install nginxingress \
  ingress-nginx/ingress-nginx --set rbac.create=true$ADDITIONAL_SET
