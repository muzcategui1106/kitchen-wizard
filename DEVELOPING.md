# Overview

This page is intended to explain how to properly develop/test the application and how to succesfully deploy it in 

# Development prerequisites

* [Docker](https://docs.docker.com/get-docker/)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
* [golang](https://go.dev/doc/install)
* [helm](https://helm.sh/docs/intro/install/)

# Starting a development environment

run `make start-development-environment`

# Deploying a local copy of the application

* run `make deploy-local VERSION=local`