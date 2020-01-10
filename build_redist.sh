#!/bin/sh

set -e


export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

mkdir -p bin/

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy -t openfaas/ofc-bootstrap:$eTAG . -f Dockerfile.redist && \
 docker create --name ofc-bootstrap openfaas/ofc-bootstrap:$eTAG && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap bin/ && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap-darwin bin/ && \
 docker cp ofc-bootstrap:/root/ofc-bootstrap.exe bin/ && \
 docker rm -f ofc-bootstrap

find bin/

