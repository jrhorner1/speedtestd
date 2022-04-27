# syntax=docker/dockerfile:1

ARG APP_VERSION="0.1.1"

### Build ###
FROM golang:alpine AS builder

RUN mkdir -p /app/bin
COPY . /app/
WORKDIR /app
RUN go build -o bin/speedtest2influx cmd/speedtest2influx/main.go 

### Deploy ###
FROM alpine:latest

ARG CLI_VERSION="1.1.1"
ARG CLI_ARCH="aarch64"
ARG CLI_URL="https://install.speedtest.net/app/cli/ookla-speedtest-${CLI_VERSION}-linux-${CLI_ARCH}.tgz" 
ARG CLI_PATH="/opt/speedtest-cli"
ARG EXPORTER_PATH=/opt/speedtest2influx

USER root

RUN mkdir -p ${EXPORTER_PATH}

COPY --from=builder /app/docs/ /app/bin/ /app/configs/ /app/LICENSE ${EXPORTER_PATH}

RUN wget -O /tmp/ookla-speedtest-${CLI_VERSION}-linux-${CLI_ARCH}.tgz "${CLI_URL}"; \
	adduser --system --uid 1000 ookla; \
	mkdir -p ${CLI_PATH}; \
	tar -xf /tmp/ookla-speedtest-${CLI_VERSION}-linux-${CLI_ARCH}.tgz -C ${CLI_PATH}; \
	chown -R 1000:1000 ${CLI_PATH}; \
	ln -sf ${CLI_PATH}/speedtest /usr/bin/speedtest; \
	mkdir -p /etc/speedtest; \
	chown -R 1000:1000 ${EXPORTER_PATH}; \
	ln -sf ${EXPORTER_PATH}/bin/speedtest2influx /usr/bin/speedtest2influx; \
	ln -sf ${EXPORTER_PATH}/configs/speedtest2influx.yaml /etc/speedtest/config.yaml; \
	apk add --no-cache libc6-compat

# Cleanup tasks
#RUN rm /tmp/ookla-speedtest-${CLI_VERSION}-linux-${CLI_ARCH}.tgz 

USER ookla

RUN speedtest --accept-license

CMD ["/usr/bin/speedtest2influx", "-config", "/etc/speedtest/config.yaml"]
