
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN /app/scripts/build.sh $TARGETOS $TARGETARCH
RUN /app/scripts/package.sh

FROM alpine:latest

ARG VERSION="0.1.0"
ARG CLI_PATH="/opt/speedtest-cli"
ARG EXPORTER_PATH=/opt/speedtest2influx

USER root

COPY --from=builder /app/build/package/ookla-speedtest-1.0.0-*-linux.tgz /app/build/package/speedtest2influx-v${VERSION}.tgz /tmp/

RUN adduser --system --uid 1000 ookla; \
	mkdir -p ${CLI_PATH}; \
	tar -xf /tmp/ookla-speedtest-1.0.0-*-linux.tgz -C ${CLI_PATH}; \
	chown -R 1000:1000 ${CLI_PATH}; \
	ln -sf ${CLI_PATH}/speedtest /usr/bin/speedtest; \
	mkdir -p ${EXPORTER_PATH}; \
	tar -xf /tmp/speedtest2influx-v${VERSION}.tgz -C ${EXPORTER_PATH}; \
	mkdir -p /etc/speedtest; \
	chown -R 1000:1000 ${EXPORTER_PATH}; \
	ln -sf ${EXPORTER_PATH}/bin/speedtest2influx /usr/bin/speedtest2influx; \
	ln -sf ${EXPORTER_PATH}/configs/speedtest2influx.yaml /etc/speedtest/config.yaml; \
	apk add --no-cache libc6-compat

USER ookla

RUN speedtest --accept-license

CMD ["/usr/bin/speedtest2influx", "-config", "/etc/speedtest/config.yaml"]