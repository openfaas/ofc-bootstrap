#!/bin/bash

helm install stable/nginx-ingress --name nginxingress --set rbac.create=true
