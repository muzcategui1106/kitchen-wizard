#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT_DIR=$SCRIPT_DIR/..
mkdir -p $ROOT_DIR/bin
cd $ROOT_DIR/cmd/kitchen-wizard 
go build -o $ROOT_DIR/bin