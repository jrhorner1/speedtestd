#!/bin/sh

sourceDir="cmd/speedtest2influx"

build_opts="-o bin/speedtest2influx"

cd /app

if [[ ! -d ./bin ]]; then
  mkdir -p ./bin
fi
go build ${build_opts} ${sourceDir}/main.go
