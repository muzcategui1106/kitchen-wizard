# Sets GIT_REF to a tag if it's present, otherwise the short git sha will be used.
GIT_REF = $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD)
# Used for Contour container image tag.
VERSION ?= $(GIT_REF)

start-development-environment:
	./scripts/start-dev-env.sh

go-build:
	mkdir -p bin
	cd cmd/kitchen-wizard; go build -o ../../bin

build:
	sudo docker build -t kitchen-wizard:$(VERSION) .

deploy-local: build
	echo $(VERSION)
	helm template development deploy/k8s/ --values deploy/k8s/values-local.yaml | sudo  kubectl apply -f /dev/stdin
	
	