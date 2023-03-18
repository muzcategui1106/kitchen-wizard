#!/bin/bash

go get -u github.com/swaggo/swag/cmd/swag
go install github.com/swaggo/swag/cmd/swag@latest
swag init
rsync -va --delete-after docs pkg/
rm -rf docs