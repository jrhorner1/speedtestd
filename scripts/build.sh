#!/bin/bash

sourceDir="cmd/speedtest2influx"

build_opts="-o bin/speedtest2influx"

if [[ ! -d ./bin ]]; then
  mkdir -p ./bin
fi
go build ${build_opts} ${sourceDir}/main.go
