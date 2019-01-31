ofc-bootstrap

"one-click" CLI to install OpenFaaS Cloud on Kubernetes

## Goals for initial release

* Install OpenFaaS and Install OpenFaaS Cloud with a single command
* Mirror features and config of OpenFaaS Cloud Community Cluster
* Use Kubernetes as the underlying provider/platform
* Use GitHub as the SCM (the source for git)
* Build via Travis
* Offer a flag for sourcing configuration from a YAML file
* Offer a dry-run flag or configuration in the YAML file
* Build a config file for the current OpenFaaS Cloud Community Cluster
* Light-weight unit-testing

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

## Stretch goals

* Automatic configuration of DNS Zones in GKE / AWS Route 53

## Non-goals

* Deep / extensive / complicated unit-tests
* Create a Docker image / run in Docker
* Installing, configuring or provisioning Kubernetes clusters or nodes
* Running on a system without bash
* Terraform/Ansible/Puppet style of experience
* Re-run without clean-up (i.e. no updates to config)
* go modules (`dep` is fine, let's add features instead)

## Pre-reqs

* Kubernetes - [development options](https://blog.alexellis.io/be-kind-to-yourself/)
    * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-binary-using-curl)
* Linux or Mac. Windows if `bash` is available
* [Go 1.10 or newer](https://golang.org/dl/)
* [dep](https://github.com/golang/dep)
* [helm](https://docs.helm.sh/using_helm/#installing-helm)
* [faas-cli](https://github.com/openfaas/faas-cli) `curl -sL https://cli.openfaas.com | sudo sh`
* OpenSSL - the `openssl` binary must be available in `PATH`

If you are using a cluster with GKE then you must run the following command:

```bash
kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
    --clusterrole=cluster-admin \
    --user="$(gcloud config get-value core/account)"
```

## Getting started

You can run ofc-bootsrap against a remote Kubernetes cluster. The instructions below use `kind` or Kubernetes in Docker to test out the instructions on your local Docker daemon.


### Get the code:

```bash
git clone https://github.com/alexellis/ofc-bootstrap
mkdir -p $GOPATH/src/github.com/openfaas-incubator
mv ofc-bootstrap $GOPATH/src/github.com/openfaas-incubator/
```

* Create a temporary cluster with `kind` or similar

```bash
go get sigs.k8s.io/kind
kind create cluster --name 1

export KUBECONFIG=$(kind get kubeconfig-path --name 1)
```

### Update `init.yaml`

* Open the Docker for Mac/Windows settings and uncheck "store my password securely" / "in a keychain"
* Run `docker login` to populate `~/.docker/config.json` - this will be used to configure your Docker registry or Docker Hub account for functions.
* Create a GitHub App and download the private key file
* If your username is not part of the default CUSTOMERS file for OpenFaaS then you should point to your own plaintext file - make sure you use the GitHub Raw CDN URL for this
* Update `init.yaml` where you see the `### User-input` section including your GitHub App's ID and the path to its private key

#### Use authz (optional)

If you'd like to restrict who can log in to just those who use a GitHub account then create a GitHub OAuth App.

Enable `auth` and fill out the required fields such as `client_secret` and `client_id`

#### Use TLS (optional)

We can automatically provision TLS certificates for your OpenFaaS Cloud cluster using the DNS01 challenge.

Pick between the following two providers for the DNS01 challenge:

* Google Cloud DNS
* AWS Route53
* DigitalOcean DNS via cert-manager 0.6.0

Configure or comment out as required in the relevant section.

You should also set up the corresponding DNS A records.

In order to enable TLS, edit the following configuration:

* Set `tls: true`
* Choose between `issuer_type: "prod"` or `issuer_type: "staging"`
* Choose between DNS Service `route53`, `clouddns` or `digitalocean` and then update init.yaml
* Go to `# DNS Service Account secret` and choose and uncomment the section you need

You can start out by using the Staging issuer, then switch to the production issuer.

* Set `issuer_type: "staging"`
* Run ofc-bootstrap with the instructions bellow

When you want to switch to the Production issuer do the following:

Flush out the staging certificates and orders

```sh
kubectl delete certificates --all  -n openfaas
kubectl delete secret -n openfaas -l="certmanager.k8s.io/certificate-name"
kubectl delete order -n openfaas --all
```

Now update the staging references to "prod":

```sh
sed -i '' s/letsencrypt-staging/letsencrypt-prod/g ./tmp/generated-ingress-ingress-wildcard.yaml
sed -i '' s/letsencrypt-staging/letsencrypt-prod/g ./tmp/generated-ingress-ingress.yaml
sed -i '' s/letsencrypt-staging/letsencrypt-prod/g ./tmp/generated-tls-auth-domain-cert.yml
sed -i '' s/letsencrypt-staging/letsencrypt-prod/g ./tmp/generated-tls-wildcard-domain-cert.yml
```

Now create the new ingress and certificates:

```sh
kubectl apply -f ./tmp/generated-ingress-ingress-wildcard.yaml
kubectl apply -f ./tmp/generated-ingress-ingress.yaml
kubectl apply -f ./tmp/generated-tls-auth-domain-cert.yml
kubectl apply -f ./tmp/generated-tls-wildcard-domain-cert.yml
```

### Run the Bootstrapper

```bash
cd $GOPATH/src/github.com/openfaas-incubator/

go build

./ofc-bootstrap -yaml=init.yaml
```

### Finish the configuration

#### Configure DNS

If you are running against a remote Kubernetes cluster you can now update your DNS entries so that they point at the IP address of your LoadBalancer found via `kubectl get svc`.

When ofc-bootstrap has completed and you know the IP of your LoadBalancer:

* `system.domain`
* `auth.system.domain`
* `*.domain`

#### Configure the GitHub App webhook

Now over on GitHub enter the URL for webhooks:

```
http://system.domain.com/github-event
```

### Smoke-test

We'll now run a smoke-test to check the dashboard shows correctly and that you can trigger a successful build.

#### View your dashboard

Now view your dashboard over at:

```
http://system.domain.com/dashboard/<username>
```

Just replace `<username>` with your GitHub account. 

#### Trigger a build

Now you can install your GitHub app on a repo, run `faas-cli new` and then rename the YAML file to `stack.yml` and do a `git push`. Your OpenFaaS Cloud cluster will build and deploy the functions found in that GitHub repo.

### Rinse & repeat

You can now edit the code or settings - do a reset of the kind or remote cluster and try again.

* Reset via kind

```
./scripts/reset-kind.sh
```

* Reset a remote cluster

```
./scripts/reset.sh
```

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

