# Sets GIT_REF to a tag if it's present, otherwise the short git sha will be used.
GIT_REF = $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD)
# Used for Contour container image tag.
VERSION ?= $(GIT_REF)

build:
	mkdir -p bin
	sudo docker build -t kitchen-wizard:$(VERSION) .

build-local:
	mkdir -p bin
	cd cmd/kitchen-wizard; go build -o ../../bin