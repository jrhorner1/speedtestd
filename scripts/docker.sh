#!/bin/sh

repo="jrhorner/ookla-speedtest"
version="0.1.0"

docker build . \
	-f ./build/package/Dockerfile \
	--tag ${repo}:${version} \
	--tag ${repo}:latest

docker push ${repo}:{$version}
docker push ${repo}:latest
