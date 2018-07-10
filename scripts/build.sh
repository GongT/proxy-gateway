#!/usr/bin/env bash

set -e

echo "building..."

protoc api/bridge.proto --go_out=plugins=grpc:./
echo " - proto"

mkdir -p dist
function build(){
	go build -o dist/$1 cmd/$1/*.go
	echo " - $1"
}

build proxy-server
build proxy-client
build single-client
