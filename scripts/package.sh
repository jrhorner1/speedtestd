#!/bin/sh

version="0.1.0"
exec=speedtest2influx
files="docs/ bin/ configs/ LICENSE"

cd /app/

tar cf build/package/speedtest2influx-v$version.tgz $files

pkgver=1.0.0
arch=$(uname -m) # Supported linux architectures: [ "aarch64", "arm", "armhf", "i386", "x86_64" ]

wget -O build/package/ookla-speedtest-1.0.0-$(uname -m)-linux.tgz "https://ookla.bintray.com/download/ookla-speedtest-1.0.0-$(uname -m)-linux.tgz"