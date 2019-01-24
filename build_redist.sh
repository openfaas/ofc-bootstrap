#!/bin/sh

export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy -t openfaas/ofc-bootstrap:$eTAG . -f Dockerfile.redist && \
 docker create --name ofc-bootstrap openfaas/ofc-bootstrap:$eTAG && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap . && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap-darwin . && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap.exe . && \
 docker rm -f ofc-bootstrap