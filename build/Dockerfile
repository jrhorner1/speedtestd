# syntax=docker/dockerfile:1
FROM alpine:latest

ARG TARGETARCH

COPY ./bin/${TARGETARCH}/speedtest /usr/bin/speedtest
COPY ./bin/${TARGETARCH}/speedtestd /usr/bin/speedtestd
COPY ./configs/speedtestd.yaml /etc/speedtestd_config.yaml

RUN speedtest --accept-license

CMD ["speedtestd", "-config", "/etc/speedtestd_config.yaml"]
