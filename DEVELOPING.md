# Overview

This page is intended to explain how to properly develop/test the application and how to succesfully deploy it in 

# Development prerequisites

* [Docker](https://docs.docker.com/get-docker/)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
* [psql](https://www.postgresql.org/docs/current/app-psql.html)
* [golang](https://go.dev/doc/install)
* [helm](https://helm.sh/docs/intro/install/)
* [golang](https://go.dev/doc/install)
* complete [local_setup](./local_setup.md)
* [npm](https://docs.npmjs.com/)

***Note*** In macOS systems you need to isntall docker-desktop as opposed to docker which includes the docker daemon. You also need to disable buildkit as it will make your builds to fail. Simply disable by going into settings. 


# getting your local IDE to be able to run the tests

* run `make grpc-build`
* go  mod download

***NOTE*** any modifications under the `pkg/pb` directory will require you run the above steps

# Starting a development environment

* run `make start-development-environment`


# Deploying a local copy of the application

***NOTE   As a one time setup you need to add the following entried to /etc/hosts for this setup to work properly***
```
<IP OF YOUR PRIMARY INTERFACE>     dex.dex.local.uzcatm-skylab.com
127.0.0.1       api.kitchen-wizard.local.uzcatm-skylab.com  ui.kitchen-wizard.local.uzcatm-skylab.com
```

* run `make deploy-local VERSION=local`

# Running in localhost (non kubernetees run)

You may need to test something quickly that might not require you to deploy to your local kubernetes cluster. If this is the case simply do the following

* on a random terminal leave the following command running
     ```
     export PGMASTER=$(kubectl get pods -o jsonpath={.items..metadata.name} -l application=spilo,cluster-name=acid-minimal-cluster,spilo-role=master -n default)
     kubectl port-forward $PGMASTER 6432:5432 -n default
     ```
* run `make run-localhost`