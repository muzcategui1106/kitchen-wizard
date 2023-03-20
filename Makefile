# Sets GIT_REF to a tag if it's present, otherwise the short git sha will be used.
GIT_REF = $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD)
# Used for Contour container image tag.
VERSION ?= $(GIT_REF)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
BUF_VERSION:=v1.15.0
SWAGGER_UI_VERSION:=v4.15.5
LOCAL_DB_PASSWORD = $(shell kubectl get secret kitchenwizard.acid-minimal-cluster.credentials.postgresql.acid.zalan.do -o 'jsonpath={.data.password}' | base64 -d)

start-development-environment:
	./scripts/start-dev-env.sh

swagger-doc-generation:
	./scripts/swagger-docs-generate.sh 

go-build: swagger-doc-generation
	mkdir -p bin
	go build -o ./bin/api

run-localhost: go-build
	./bin/api --dex-provider-url "https://dex.dex.local.uzcatm-skylab.com" \
	--oidc-client-id example-app \
	--oidc-client-secret ZXhhbXBsZS1hcHAtc2VjcmV0 \
	--oidc-redirect-url "http://localhost:8443" \
	--postgres-db-hostname localhost \
	--postgres-db-username kitchenwizard \
	--postgres-db-port "6432" \
	--postgres-db-password $(LOCAL_DB_PASSWORD)

build-local:
	echo "building image"
	docker build -t kitchen-wizard:local .

deploy-local: build-local
	echo "loading image to local cluster"
	kind load docker-image kitchen-wizard:local

	echo "generating k8s manifests"
	DB_PASSWORD=`$(kubectl get secret postgres.acid-minimal-cluster.credentials.postgresql.acid.zalan.do -o 'jsonpath={.data.password}' | base64 -d)`
	helm template development deploy/k8s/ --values deploy/k8s/values-local.yaml --set dbPassword=$(LOCAL_DB_PASSWORD) | kubectl apply -f /dev/stdin
	
	echo "sleeping 5 seconds to ensure image has gotten to nodes"
	sleep 5

	echo "restarting deployments to ensure latest version is used"
	kubectl -n kitchen-wizard rollout restart deployments
	
	