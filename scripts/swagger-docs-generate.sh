#!/bin/bash

./get-swag-dependencies.sh
swag init
rsync -va --delete-after docs pkg/
rm -rf docs