#!/bin/sh

version="0.1.0"

docker buildx build . \
	--platform linux/amd64,linux/arm64,linux/arm \
	-f ./build/package/Dockerfile \
	--tag speedtest2influx:${version} \
	--tag speedtest2influx:latest #\
#	--load
	
