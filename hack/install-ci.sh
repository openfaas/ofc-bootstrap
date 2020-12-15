#!/bin/bash

set -e

curl -sLSf https://get.docker.com | sudo sh

echo "export GOPATH=\$HOME/go" | tee -a ~/.bash_profile
echo "export PATH=\$GOPATH/bin:\$PATH:/usr/local/go/bin/" | tee -a ~/.bash_profile

curl -sLSf https://dl.get-arkade.dev | sudo sh

arkade get kind
sudo mv $HOME/.arkade/bin/kind /usr/local/bin/

./bin/ofc-bootstrap version
