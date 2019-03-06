## Roadmap

### Goals for initial release

* Install OpenFaaS and Install OpenFaaS Cloud with a single command
* Mirror features and config of OpenFaaS Cloud Community Cluster
* Use Kubernetes as the underlying provider/platform
* Use GitHub as the SCM (the source for git)
* Build via Travis
* Offer a flag for sourcing configuration from a YAML file
* Offer a dry-run flag or configuration in the YAML file
* Build a config file for the current OpenFaaS Cloud Community Cluster
* Light-weight unit-testing

### Goals for 1.0

* Publish a static binary on GitHub Releases for `ofc-bootstrap` tool
* Use GitLab for as SCM (the source for git)
* Implement a back-end for Swarm in addition to the Kubernetes support.
* Allow namespaces to be overriden from `openfaas`/`openfaas-fn` to something else

### Goals for 2.0

* Build a suitable dev environment for local work (without Ingress, TLS)
* Add version number to YAML file i.e `1.0` to enable versioning/migration of configs
* Move code into official CLI via `faas-cli system install openfaas-cloud`
* Separate out the OpenFaaS installation for the official CLI `faas-cli system install --kubernetes/--swarm`

### Stretch goals

* Automatic configuration of DNS Zones in GKE / AWS Route 53

### Non-goals

* Deep / extensive / complicated unit-tests
* Create a Docker image / run in Docker
* Installing, configuring or provisioning Kubernetes clusters or nodes
* Running on a system without bash
* Terraform/Ansible/Puppet style of experience
* Re-run without clean-up (i.e. no updates to config)
* go modules (`dep` is fine, let's add features instead)



## Status

Help is wanted - the code is in a private repo for OpenFaaS maintainers to contribute to. Sign-off/DCO is required and standard OpenFaaS contributing procedures apply.

Status:
* [ ] Flag: Add dry-run to init.yaml
* [x] Step: generate `payload_secret` for trust
* [x] Refactor: default to init.yaml if present
* [x] Step: Clone OpenFaaS Cloud repo https://github.com/openfaas/openfaas-cloud
* [x] Step: deploy container builder (buildkit)
* [x] Step: Add Ingress controller
* [x] Step: Install OpenFaaS via helm
* [x] Step: Install tiller sa
* [x] Step: Install OpenFaaS namespaces
* [x] Wildcard ingress
* [x] Auth ingress
* [x] init.yml - define GitHub App and load via struct
* [x] Step: deploy OpenFaaS Cloud primary functions
* [x] Step: deploy OpenFaaS Cloud dashboard
* [x] Template: dashboard stack.yml if required
* [x] Template: `gateway_config.yml`
* [x] Step: install SealedSecrets
* [x] Step: export SealedSecrets pub-cert
* [ ] Step: export all passwords required for user such as GW via `kubectl`
* [x] Step: setup issuer and certificate entries for cert-manager (probably with staging cert?) - make this optional to prevent rate-limiting.
* [x] Make TLS optional in the Ingress config (not to get rate-limited by LetsEncrypt)
* [x] init.yml - add `github_app_id` and `WEBHOOK_SECRET`
* [x] Create basic-auth secrets for the functions in `openfaas-fn`
* [x] Step: Install Minio and generate keys
* [x] init.yml - define and OAuth App and load via struct
* [x] Step: generate secrets and keys for the auth service (see auth/README.md)
* [x] Template: auth service deployment YAML file
* [ ] Refactor: Generate passwords via Golang code or library

Add all remaining steps from [installation guide](https://github.com/openfaas/openfaas-cloud/tree/master/docs).

Caveats:

* JetStack's cert-manager is currently pinned to an earlier version due to issues with re-creating the CRD entries. 
