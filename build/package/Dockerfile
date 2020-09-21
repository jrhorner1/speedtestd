FROM apline:latest

ARG PKGVER = "1.0.0"
# Supported linux architectures: [ "aarch64", "arm", "armhf", "i386", "x86_64" ]
ARG ARCH = "arm"
ARG CLI_URL = "https://ookla.bintray.com/download/ookla-speedtest-${PKGVER}-${ARCH}-linux.tgz"
ARG CLI_PATH = "/opt/speedtest-cli"

ARG EXPORTER_PATH = /opt/speedtest2influx

ENV INFLUXDB_ADDRESS
ENV INFLUXDB_PORT 
ENV INFLUXDB_DATABASE
ENV INFLUXDB_USERNAME
ENV INFLUXDB_PASSWORD
ENV SPEEDTEST_SERVER_ID
ENV SPEEDTEST_SERVER_NAME
ENV SPEEDTEST_INTERVAL
ENV LOGGING_LEVEL

USER root

RUN groupadd -r -g 1000 ookla \
    useradd -r -u 1000 -g 1000 ookla

ADD --chown=1000:1000 ${CLI_URL} ${CLI_PATH}
RUN tar -xf ${CLI_PATH}/ookla-speedtest-1.0.0-${ARCH}-linux.tgz -C ${CLI_PATH} \
    chown -R 1000:1000 ${CLI_PATH} \
    ln -sf ${CLI_PATH}/speedtest /usr/bin/speedtest

RUN mkdir -p ${EXPORTER_PATH} /etc/speedtest
COPY bin/ config/ LICENSE README.md ${EXPORTER_PATH}
RUN chown -R 1000:1000 ${EXPORTER_PATH} \
    ln -sf ${EXPORTER_PATH}/bin/speedtest2influx /usr/bin/speedtest2influx \
    ln -sf ${EXPORTER_PATH}/config/config.yaml /etc/speedtest/config.yaml

USER ookla

CMD ["/usr/bin/speedtest2influx", "-config", "/etc/speedtest/config.yaml"]