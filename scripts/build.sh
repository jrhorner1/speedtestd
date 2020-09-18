#!/bin/bash

sourceDir="cmd/speedtest2influx"

build_opts="-o speedtest2influx"

go build ${build_opts} ${sourceDir}/main.go