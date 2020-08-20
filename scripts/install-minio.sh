#!/bin/bash

export ACCESS_KEY=$(kubectl get secret -n openfaas-fn s3-access-key -o jsonpath='{.data.s3-access-key}'| base64 --decode)
export SECRET_KEY=$(kubectl get secret -n openfaas-fn s3-secret-key -o jsonpath='{.data.s3-secret-key}'| base64 --decode)

# Minio has a default requests value of 4Gi RAM
# https://github.com/minio/charts/blob/master/minio/values.yaml

helm upgrade --install cloud-minio --namespace openfaas \
   --set accessKey="$ACCESS_KEY",secretKey="$SECRET_KEY",replicas=1,persistence.enabled=false,service.port=9000,service.type=NodePort \
   --set resources.requests.memory=512Mi \
  minio/minio

