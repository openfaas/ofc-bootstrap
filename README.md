## ofc-bootstrap

> Provide a managed OpenFaaS experience for your team

How? By automating the whole installation of OpenFaaS Cloud on Kubernetes into a single command and CLI.

[![Build Status](https://travis-ci.com/openfaas-incubator/ofc-bootstrap.svg?branch=master)](https://travis-ci.com/openfaas-incubator/ofc-bootstrap)

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

### Conceptual design

The ofc-bootstrap tool is used to install OpenFaaS Cloud in a single click. You will need to configure it with all the necessary secrets and configuration beforehand using a YAML file.

![](./docs/ofc-bootstrap.png)

> ofc-bootstrap packages a number of primitives such as an IngressController, a way to obtain certificates from LetsEncrypt, the OpenFaaS Cloud components, OpenFaaS itself and Minio for build log storage. Each component is interchangeable.

### Video demo

View a video demo by Alex Ellis running `ofc-bootstrap` in around 100 seconds on DigitalOcean.

[![View demo](https://img.youtube.com/vi/Sa1VBSfVpK0/0.jpg)](https://www.youtube.com/watch?v=Sa1VBSfVpK0)

## Roadmap

See the [ROADMAP.md](./ROADMAP.md) for features, development status and backlogs. 

## Get started

Follow the [user guide](USER_GUIDE.md).

### Join us on Slack

Got questions, comments or suggestions?

Join the team and community over on [Slack](https://docs.openfaas.com/community)

