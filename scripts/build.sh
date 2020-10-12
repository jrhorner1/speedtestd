#!/bin/sh

sourceDir="cmd/speedtest2influx"

build_opts="-o bin/speedtest2influx"

if [[ ! -d ./bin ]]; then
  mkdir -p ./bin
fi
#GO_OS=$1 GO_ARCH=$2 
go build ${build_opts} ${sourceDir}/main.go

