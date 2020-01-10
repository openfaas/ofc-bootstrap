#!/bin/sh

set -e

export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

echo Building openfaas/ofc-bootstrap:$eTAG
mkdir -p bin

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy -t openfaas/ofc-bootstrap:$eTAG . && \
 docker create --name ofc-bootstrap openfaas/ofc-bootstrap:$eTAG && \
 docker cp ofc-bootstrap:/usr/bin/ofc-bootstrap bin/. && \
 docker rm -f ofc-bootstrap
