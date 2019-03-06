[![Build Status](https://travis-ci.org/openfaas-incubator/ofc-bootstrap.svg?branch=master)](https://travis-ci.org/openfaas-incubator/ofc-bootstrap)
ofc-bootstrap

### What is this and who is it for?

> "one-click" CLI to install OpenFaaS Cloud on Kubernetes

You can use this tool to configure a Kubernetes cluster with OpenFaaS Cloud. You just need to complete all the pre-requisites and fill out your `init.yaml` file then run the tool. It automates several pages of manual steps using Golang templates and bash scripts so that you can get your own OpenFaaS Cloud in around 1.5 minutes.

Experience level: intermediate Kubernetes/cloud.

The `ofc-bootstrap` will install the following components:

* OpenFaaS installed with helm
* Nginx as your IngressController - with rate-limits configured
* SealedSecrets from Bitnami - store secrets for functions in git
* cert-manager - provision HTTPS certificates with LetsEncrypt
* Docker’s buildkit - to building immutable Docker images for each function
* Authentication/authorization - through OAuth2 delegating to GitHub/GitLab
* Deep integration into GitHub/GitLab - for updates and commit statuses
* A personalized dashboard for each user

## Roadmap

See the [ROADMAP.md](./ROADMAP.md) for features, development status and backlogs. 

## Get started

To run a production-quality OpenFaaS Cloud then execute `ofc-bootstrap` with a `kubeconfig` file pointing to a remote Kubernetes service. For development and testing you can use the instructions below with `kind`. The `kind` distribution of Kubernetes does not require anything on your host other than Docker.

### Pre-reqs

This tool automates the installation of OpenFaaS Cloud on Kubernetes. It can be set up on a public cloud provider with a managed Kubernetes offering, where a `LoadBalancer` is available. If you are deploying to a cloud or Kubernetes cluster where the type `LoadBalancer` is unavailable then you will need to change `ingress: loadbalancer` to `ingress: host` in `init.yaml`. This will provision Nginx as a `DaemonSet` exposed on port `80` and `443`.

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

First run `cp example.init.yaml init.yaml` to get your own example `init.yaml` file.

* Open the Docker for Mac/Windows settings and uncheck "store my password securely" / "in a keychain"
* Run `docker login` to populate `~/.docker/config.json` - this will be used to configure your Docker registry or Docker Hub account for functions.
* Create a GitHub App and download the private key file
  * Read the docs for how to [configure your GitHub App](https://docs.openfaas.com/openfaas-cloud/self-hosted/github/)
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

Then you need to enter the Webhook secret that was generated during the bootstrap process. Run the following commands to extract and decode it:

> If you have `jq` installed, this one-liner would be handy: `kubectl -n openfaas-fn get secret github-webhook-secret -o json | jq '.data | map_values(@base64d)'`. Otherwise, continue below.

```
$ kubectl -n openfaas-fn get secret github-webhook-secret -o yaml
```

This spits out the Secret object definition, including a field like:

```
data:
  github-webhook-secret: <redacted base64-encoded string>
```

Copy out that string and decode it:

```
$ echo 'redacted base64-encoded string' | base64 --decode
```

Enter this value into the webhook secret field in your GitHub App.

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

