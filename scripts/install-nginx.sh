#!/bin/bash

echo helm install stable/nginx-ingress --name nginxingress --set rbac.create=true$ADDITIONAL_SET
helm install stable/nginx-ingress --name nginxingress --set rbac.create=true$ADDITIONAL_SET
