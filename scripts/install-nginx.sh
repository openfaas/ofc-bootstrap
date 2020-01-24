#!/bin/bash

set -e

echo helm upgrade --install nginxingress \
  stable/nginx-ingress --set rbac.create=true$ADDITIONAL_SET
helm upgrade --install nginxingress \
  stable/nginx-ingress --set rbac.create=true$ADDITIONAL_SET
