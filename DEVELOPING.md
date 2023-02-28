# Overview

This page is intended to explain how to properly develop/test the application and how to succesfully deploy it in 

# Development prerequisites

* [Docker](https://docs.docker.com/get-docker/)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
* [golang](https://go.dev/doc/install)
* [helm](https://helm.sh/docs/intro/install/)
* [golang](https://go.dev/doc/install)
* complete [local_setup](./local_setup.md)

***Note*** In macOS systems you need to isntall docker-desktop as opposed to docker which includes the docker daemon. You also need to disable buildkit as it will make your builds to fail. Simply disable by going into settings. 

# getting your local IDE to be able to run the tests

* run `make grpc-build`
* go  mod download

***NOTE*** any modifications under the `pkg/pb` directory will require you run the above steps

# Starting a development environment

* run `make start-development-environment`


# Deploying a local copy of the application

* run `make deploy-local VERSION=local`