ofc-bootstrap

"one-click" CLI to install OpenFaaS Cloud on Kubernetes

## Goals for initial release

* Install OpenFaaS and Install OpenFaaS Cloud with a single command
* Mirror OpenFaaS Cloud Community Cluster features/coverage
* Use Kubernetes as the underlying provider/platform
* Use GitHub as the SCM (the source for git)
* Build via Travis
* Offer a flag for sourcing configuration from a YAML file
* Offer a dry-run flag or configuration in the YAML file
* Build a config file for the current OpenFaaS Cloud Community Cluster

## Goals for 1.0

* Publish a static binary on GitHub Releases for `ofc-bootstrap` tool
* Use GitLab for as SCM (the source for git)
* Implement a back-end for Swarm in addition to the Kubernetes support.
* Allow namespaces to be overriden from `openfaas`/`openfaas-fn` to something else

## Goals for 2.0

* Build a suitable dev environment for local work (without Ingress, TLS)
* Add version number to YAML file i.e `1.0` to enable versioning/migration of configs
* Move code into official CLI via `faas-cli system install openfaas-cloud`
* Separate out the OpenFaaS installation for the official CLI `faas-cli system install --kubernetes/--swarm`

## Non-goals

* Create a Docker image / run in Docker
* Installing, configuring or provisioning Kubernetes clusters or nodes
* Running on a system without bash
* Terraform/Ansible/Puppet style of experience
* Re-run without clean-up (i.e. no updates to config)
* go modules (`dep` is fine, let's add features instead)

## Pre-reqs

* Kubernetes - [development options](https://blog.alexellis.io/be-kind-to-yourself/)
* Linux or Mac. Windows if `bash` is available
* [Go 1.10 or newer](https://golang.org/dl/)
* [dep](https://github.com/golang/dep)
* [helm](https://docs.helm.sh/using_helm/#installing-helm)
* [faas-cli](https://github.com/openfaas/faas-cli) `curl -sL https://cli.openfaas.com | sudo sh`

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

* Reset the cluster

```
./scripts/reset-kind.sh
```

Now you can edit the code and run it again. Tiller takes several seconds to come up.

Notes:

JetStack's cert-manager is currently pinned to an earlier version due to issues with re-creating the CRD entries. 

## Status

Help is wanted - the code is in a private repo for OpenFaaS maintainers to contribute to. Sign-off/DCO is required and standard OpenFaaS contributing procedures apply.

Status:
* [ ] Library: Generate passwords via Golang code or library
* [ ] Step: Add Ingress controller
* [ ] Step: Install OpenFaaS via helm
* [ ] Step: Install tiller sa
* [ ] Step: Install OpenFaaS namespaces
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
* [ ] Step: deploy OpenFaaS Cloud primary functions
* [ ] Step: deploy OpenFaaS Cloud dashboard
* [ ] Step: deploy container builder (buildkit)

Add all remaining steps from [installation guide](https://github.com/openfaas/openfaas-cloud/tree/master/docs).

