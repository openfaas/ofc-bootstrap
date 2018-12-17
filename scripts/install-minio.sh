#!/bin/bash

export ACCESS_KEY=$(kubectl get secret -n openfaas-fn s3-access-key -o jsonpath='{.data.s3-access-key}'|base64 -D)
export SECRET_KEY=$(kubectl get secret -n openfaas-fn s3-secret-key -o jsonpath='{.data.s3-secret-key}'|base64 -D)

helm install --name cloud-minio --namespace openfaas \
   --set accessKey=$ACCESS_KEY,secretKey=$SECRET_KEY,replicas=1,persistence.enabled=false,service.port=9000,service.type=NodePort \
  stable/minio
