## ofc-bootstrap

> "one-click" CLI to install OpenFaaS Cloud on Kubernetes

[![Build Status](https://travis-ci.org/openfaas-incubator/ofc-bootstrap.svg?branch=master)](https://travis-ci.org/openfaas-incubator/ofc-bootstrap)

### What is this and who is it for?

You can use this tool to configure a Kubernetes cluster with [OpenFaaS Cloud](https://github.com/openfaas/openfaas-cloud). You just need to complete all the pre-requisites and fill out your `init.yaml` file then run the tool. It automates several pages of manual steps using Golang templates and bash scripts so that you can get your own [OpenFaaS Cloud](https://github.com/openfaas/openfaas-cloud) in around 1.5 minutes.

Experience level: intermediate Kubernetes & cloud.

The `ofc-bootstrap` will install the following components:

* [OpenFaaS](https://github.com/openfaas/faas) installed with helm
* [Nginx as your IngressController](https://github.com/kubernetes/ingress-nginx) - with rate-limits configured
* [SealedSecrets](https://github.com/bitnami-labs/sealed-secrets) from Bitnami - store secrets for functions in git
* [cert-manager](https://github.com/jetstack/cert-manager) - provision HTTPS certificates with LetsEncrypt
* [buildkit from Docker](https://github.com/moby/buildkit) - to building immutable Docker images for each function
* Authentication/authorization - through OAuth2 delegating to GitHub/GitLab
* Deep integration into GitHub/GitLab - for updates and commit statuses
* A personalized dashboard for each user

### Video demo

View a video demo by Alex Ellis running `ofc-bootstrap` in around 100 seconds on DigitalOcean.

[![View demo](https://img.youtube.com/vi/Sa1VBSfVpK0/0.jpg)](https://www.youtube.com/watch?v=Sa1VBSfVpK0)

## Roadmap

See the [ROADMAP.md](./ROADMAP.md) for features, development status and backlogs. 

## Installation

To run a production-quality OpenFaaS Cloud then execute `ofc-bootstrap` with a `kubeconfig` file pointing to a remote Kubernetes service. For development and testing you can use the instructions below with `kind`. The `kind` distribution of Kubernetes does not require anything on your host other than Docker.

### Pre-reqs

This tool automates the installation of OpenFaaS Cloud on Kubernetes. We will now install some required tools and then create either a local or remote cluster.

#### Tools

* Kubernetes - [development options](https://blog.alexellis.io/be-kind-to-yourself/)
    * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-binary-using-curl)
* Linux or Mac. Windows if `bash` is available
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

#### Create a cluster

* Create a production cluster

You can create a managed or self-hosted Kubernetes cluster using a Kubernetes engine such as GKE, AWS, DigitalOcean or by using `kubeadm`. Once set up make sure you have set your `KUBECONFIG` and / or `kubectl` tool to point at a the new cluster.

* Create a local cluster for testing

For testing you can create a local cluster using `kind`, `minikube` or Docker Desktop. This is how you can install `kind` to setup a local cluster in a Docker container.

First install [Go 1.10 or newer](https://golang.org/dl/)

Now use `go get` to install `kind` and point your `KUBECONFIG` variable at the new cluster.

```bash
go get sigs.k8s.io/kind
kind create cluster --name 1

export KUBECONFIG=$(kind get kubeconfig-path --name 1)
```

### Get `ofc-bootstrap`

* Clone the GitHub repo

```bash
mkdir -p $GOPATH/src/github.com/openfaas-incubator
cd $GOPATH/src/github.com/openfaas-incubator/
git clone https://github.com/openfaas-incubator/ofc-bootstrap
```

* Download the latest binary release from GitHub

Download [ofc-boostrap](https://github.com/openfaas-incubator/ofc-bootstrap/releases) from the GitHub releases page and move it to `/usr/local/bin/`. You may also need to run `chmod +x /usr/local/bin/ofc-bootstrap`.

### Create your own `init.yaml`

First run `cp example.init.yaml init.yaml` to get your own `init.yaml` file.

Log into your Docker registry or the Docker Hub:

* Open the Docker for Mac/Windows settings and uncheck "store my password securely" / "in a keychain"
* Run `docker login` to populate `~/.docker/config.json` - this will be used to configure your Docker registry or Docker Hub account for functions.

Choose SCM between GitHub and GitLab, by setting `scm: github` or `scm: gitlab`

Setup the GitHub / GitLab App and OAuth App

* For GitHub create a GitHub App and download the private key file
  * Read the docs for how to [configure your GitHub App](https://docs.openfaas.com/openfaas-cloud/self-hosted/github/)
  * Update `init.yaml` where you see the `### User-input` section including your GitHub App's ID, Webhook secret and the path to its private key
* For GitLab create a System Hook
  * Update the `### User-input` section including your System Hook's API Token and Webhook secret
* Create your GitHub / GitLab OAuth App which is used for logging in to the dashboard
* For GitLab update `init.yaml` with your `gitlab_instance`

Create your own GitHub repo with a CUSTOMERS ACL file

* Create a new public GitHub repo
* Add a file named `CUSTOMERS` and place each username or GitHub org you will use on a separate line
* Add the GitHub RAW CDN URL into the init.yaml file

#### Decide if you're using a LoadBalancer

It can be set up on a public cloud provider with a managed Kubernetes offering, where a `LoadBalancer` is available. If you are deploying to a cloud or Kubernetes cluster where the type `LoadBalancer` is unavailable then you will need to change `ingress: loadbalancer` to `ingress: host` in `init.yaml`. This will provision Nginx as a `DaemonSet` exposed on port `80` and `443`.

#### Use authz (optional)

If you'd like to restrict who can log in to just those who use a GitHub account then create a GitHub OAuth App.

Enable `auth` and fill out the OAuth App `client_id`. Configure `of-client-secret` with the OAuth App Client Secret.
For GitLab set your `oauth_provider_base_url`.

#### Use TLS (optional)

We can automatically provision TLS certificates for your OpenFaaS Cloud cluster using the DNS01 challenge.

Pick between the following providers for the DNS01 challenge:

* Google Cloud DNS
* AWS Route53
* DigitalOcean DNS via cert-manager 0.6.0

> Note: At time of writing DigitalOcean are offering free management of DNS.

Configure or comment out as required in the relevant section.

You should also set up the corresponding DNS A records in your DNS management dashboard after finishing all the steps in this guide.

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

### Run the `ofc-bootstrap`

```bash
cd $GOPATH/src/github.com/openfaas-incubator/ofc-bootstrap

./ofc-bootstrap -yaml=init.yaml
```

### Finish the configuration

If you get anything wrong, don't worry you can use the `./scripts/reset.sh` file to remove all the components. Then edit `init.yaml` and start over. Be careful running this script and make 100% sure that you are pointing at the correct cluster. 

#### Configure DNS

If you are running against a remote Kubernetes cluster you can now update your DNS entries so that they point at the IP address of your LoadBalancer found via `kubectl get svc`.

When ofc-bootstrap has completed and you know the IP of your LoadBalancer:

* `system.domain`
* `auth.system.domain`
* `*.domain`

#### Configure the GitHub / GitLab App Webhook

Now over on GitHub / GitLab enter the URL for webhooks:

GitHub:
```
http://system.domain.com/github-event
```
GitLab:
```
http://system.domain.com/gitlab-event
```

For more details see the [GitLab instructions](https://github.com/openfaas/openfaas-cloud/blob/master/docs/GITLAB.md) in OpenFaaS Cloud.

Then you need to enter the Webhook secret that was generated during the bootstrap process. Run the following commands to extract and decode it:

```echo $(kubectl get secret -n openfaas-fn github-webhook-secret -o jsonpath="{.data.github-webhook-secret}" | base64 --decode; echo)```

Open the Github App UI and paste in the value into the "Webhook Secret" field.

### Smoke-test

Now run a smoke-test to check the dashboard shows correctly and that you can trigger a successful build.

#### View your dashboard

Now view your dashboard over at:

```
http://system.domain.com/dashboard/<username>
```

Just replace `<username>` with your GitHub account. 

#### Trigger a build

Now you can install your GitHub app on a repo, run `faas-cli new` and then rename the YAML file to `stack.yml` and do a `git push`. Your OpenFaaS Cloud cluster will build and deploy the functions found in that GitHub repo.

#### Something went wrong?

If you think that everything is set up correctly but want to troubleshoot then head over to the GitHub App webpage and click "Advanced" - here you can find each request/response from the GitHub push events. You can resend them or view any errors.

#### Invite your team

For each user or org you want to enroll into your OpenFaaS Cloud edit the CUSTOMERS ACL file and add their username on a new line. For example if I wanted the user `alexellis` and the org `openfaas` to host git repos containing functions:

```
openfaas
alexellis
```

### Join us on Slack

Got questions, comments or suggestions?

Join the team and community over on [Slack](https://docs.openfaas.com/community)