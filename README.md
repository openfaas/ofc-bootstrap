ofc-bootstrap

"one-click" CLI to install OpenFaaS Cloud on Kubernetes

## Goals

* Get OpenFaaS and OpenFaaS Cloud installed with a single command
* Mirror OpenFaaS Cloud Community Cluster features/coverage
* Build an environment for Kubernetes

## Stretch goals

* Library for re-use in official CLI or other dedicated CLI i.e. `faas-cli system install openfaas-cloud`
* Command for installing a base-OpenFaaS system `faas-cli system install --kubernetes/--swarm` 
* Build a suitable dev environment for local work
* Build environment out for a Swarm cluster

## Non-goals

* Running in a Docker image
* Installing, configuring or provisioning Kubernetes clusters or nodes
* Running on a system without bash
* Terraform/Ansible/Puppet style of experience
* Re-run without clean-up
* Updates to config after initial deployment
* To use go modules

## Pre-reqs

* [Go 1.10 or newer](https://golang.org/dl/)
* [dep](https://github.com/golang/dep)
* [helm](https://docs.helm.sh/using_helm/#installing-helm)
* Kubernetes
* Linux, Mac or Windows with bash available

## Getting started

* Get the code:

```bash
git clone https://github.com/alexellis/ofc-bootstrap
mkdir -p $GOPATH/src/github.com/alexellis
mv ofc-bootstrap $GOPATH/src/github.com/alexellis/
```

* Create a temporary cluster with `kind` or similar

```bash
go get sigs.k8s.io/kind
kind create cluster
export KUBECONFIG=$(kind get kubeconfig-path)
```

* Run the code

```bash
cd $GOPATH/src/github.com/alexellis/
mkdir -p tmp
go run main.go -yaml init.yaml
```

## Status

Help is wanted - the code is in a private repo for OpenFaaS maintainers to contribute to. Sign-off/DCO is required and standard OpenFaaS contributing procedures apply.

[ ] Dry-run flag needed
[ ] Exec commands need to be actioned, but are just printed
[ ] Secret generation isn't working, but should be moved to Golang code - perhaps using a popular Go library?

