ofc-bootstrap

"one-click" CLI to install OpenFaaS Cloud on Kubernetes

## Goals

* Install OpenFaaS and Install OpenFaaS Cloud with a single command
* Mirror OpenFaaS Cloud Community Cluster features/coverage
* Use Kubernetes as the underlying provider/platform
* Use GitHub as the SCM (the source for git)
* Offer a flag for sourcing configuration from a YAML file
* Offer a dry-run flag or configuration in the YAML file
* Build a config file for the current OpenFaaS Cloud Community Cluster

## Stretch goals

* Publish a static binary
* Use GitLab for as SCM (the source for git)
* Build environment out for a Swarm cluster
* Add version to YAML file to enable versioning/migration of configs
* Create a Go library for re-use in official CLI or other dedicated CLI i.e. `faas-cli system install openfaas-cloud`
* Build a suitable dev environment for local work
* Add command for installing a base-OpenFaaS system `faas-cli system install --kubernetes/--swarm` 

## Non-goals

* Create a Docker image / run in Docker
* Installing, configuring or provisioning Kubernetes clusters or nodes
* Running on a system without bash
* Terraform/Ansible/Puppet style of experience
* Re-run without clean-up (i.e. no updates to config)
* go modules (`dep` is fine, let's add features instead)

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

Status:
* [ ] Flag: Add dry-run
* [ ] Exec commands need to be actioned, but are just printed
* [ ] Secret generation isn't working, but should be moved to Golang code - perhaps using a popular Go library?
* [x] Wildcard ingress
* [x] Auth ingress
* [ ] Template: auth service YAML
* [ ] Template: dashboard stack.yml
* [x] Template: `gateway_config.yml`
* [ ] Step: install SealedSecrets
* [ ] Step: export SealedSecrets pub-cert
* [ ] Step: export all passwords required for user such as GW via `kubectl`
* [ ] Step: setup issuer and certificate entries for cert-manager
* [ ] Step: generate `payload_secret` for trust
* [ ] init.yml - add `github_app_id` and `WEBHOOK_SECRET`
* [ ] Create basic-auth secrets for the functions in `openfaas-fn`
* [ ] Step: Install Minio and generate keys
* [ ] init.yml - define GitHub App and OAuth App and load via struct


Add all remaining steps from [installation guide](https://github.com/openfaas/openfaas-cloud/tree/master/docs).



