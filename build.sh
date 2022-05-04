#!/bin/bash

# Speedtest CLI variables, ref. https://www.speedtest.net/apps/cli
VERSION="1.1.1"
OS="linux"
ARCHITECTURES=( x86_64 aarch64 armhf )

go_build () { # BEGIN go build function
    mkdir -p bin/$1 
    printf "Compiling speedtestd for $1\n\n"
    GOOS=linux GOARCH=$1 go build -o bin/$1/speedtestd speedtestd/*.go
} # END

unpack () { # BEGIN speedtest unpack function
    stat bin/$1/speedtest > /dev/null 
    if [[ ! $? -eq 0 ]]; then
        mkdir -p bin/$1
        ARCH=$2
        printf "Downloading Ookla Speedtest CLI\n\n"
        curl -OL "https://install.speedtest.net/app/cli/ookla-speedtest-${VERSION}-${OS}-${ARCH}.tgz"
        printf "Unpacking Ookla Speedtest CLI\n\n"
        tar xf ookla-speedtest-${VERSION}-${OS}-${ARCH}.tgz --include "speedtest" -C bin/$1
        rm ookla-speedtest-${VERSION}-${OS}-${ARCH}.tgz
    fi
} # END

# Loop through each architecture, download the cli, build the app, then cleanup 
for ARCH in "${ARCHITECTURES[@]}"; do
    case ${ARCH} in
        "x86_64") 
            go_build "amd64" 
            unpack "amd64" ${ARCH}
            ;; 
        "aarch64") 
            go_build "arm64" 
            unpack "arm64" ${ARCH} 
            ;;
        "armhf") 
            go_build "arm" 
            unpack "arm" ${ARCH}
            ;; 
    esac
done

# Application variables
REGISTRY="docker.io"
USER="jrhorner"
REPOSITORY="speedtestd"
VERSION="0.2.1"

# Build the multiarch container
printf "Building multi-arch docker container\n\n"
docker buildx build . \
    --platform "linux/amd64,linux/arm64/v8,linux/arm/v7" \
    --tag ${REGISTRY}/${USER}/${REPOSITORY}:${VERSION} \
    --tag ${REGISTRY}/${USER}/${REPOSITORY}:latest \
    --push
