#!/bin/bash

SCRIPT_PATH=`readlink -f "$0"`
SCRIPT_DIR=`dirname "$SCRIPT_PATH"`

$SCRIPT_DIR/get-swag-dependencies.sh
swag init
rsync -va --delete-after docs pkg/
rm -rf docs