# Sets GIT_REF to a tag if it's present, otherwise the short git sha will be used.
GIT_REF = $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD)
# Used for Contour container image tag.
VERSION ?= $(GIT_REF)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
BUF_VERSION:=v1.9.0
SWAGGER_UI_VERSION:=v4.15.5

start-development-environment:
	./scripts/start-dev-env.sh


grpc-build:
	#go get github.com/grpc-ecosystem/grpc-gateway
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) mod update
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate

go-build:
	mkdir -p bin
	cd cmd/api; go build -o ../../bin

build:
	echo "building image"
	sudo docker build -t kitchen-wizard:$(VERSION) .

deploy-local: build
	echo "loading image to local cluster"
	sudo kind load docker-image kitchen-wizard:$(VERSION)

	echo "generating k8s manifests"
	helm template development deploy/k8s/ --values deploy/k8s/values-local.yaml | sudo  kubectl apply -f /dev/stdin
	

	echo "sleeping 5 seconds to ensure image has gotten to nodes"
	sleep 5

	echo "restarting deployments to ensure latest version is used"
	sudo kubectl -n kitchen-wizard rollout restart deployments
	
	