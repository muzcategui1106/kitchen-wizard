# Sets GIT_REF to a tag if it's present, otherwise the short git sha will be used.
GIT_REF = $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD)
# Used for Contour container image tag.
VERSION ?= $(GIT_REF)

start-development-environment:
	./scripts/start-dev-env.sh

grpc-build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/server/kitchen_wizard.proto 

go-build:
	mkdir -p bin
	cd cmd/kitchen-wizard; go build -o ../../bin

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
	
	