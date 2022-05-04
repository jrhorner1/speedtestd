# syntax=docker/dockerfile:1
FROM alpine:latest

ARG TARGETARCH

COPY ./bin/${TARGETARCH}/speedtest /usr/bin/speedtest
COPY ./bin/${TARGETARCH}/speedtestd /usr/bin/speedtestd
COPY ./config.yaml /etc/speedtestd.yaml

RUN speedtest --accept-license

CMD ["speedtestd", "-c", "/etc/speedtestd.yaml"]
