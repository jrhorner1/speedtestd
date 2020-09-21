#!/bin/bash

version="0.1.0"

docker build build/package/ -t speedtest2influx:latest -t speedtest2influx:${version}