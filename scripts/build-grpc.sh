#!/bin/sh -x

env

protoc -I/usr/local/include -I. \
	-I$(INCLUDE_GOOGLE_APIS) \
    -I$(INCLUDE_GRPC_GATEWAY) \
    --go_out=. --go_opt=paths=source_relative --go-grpc_out=/gen/ --go-grpc_opt=paths=source_relative /gen/kitchen_wizard.proto