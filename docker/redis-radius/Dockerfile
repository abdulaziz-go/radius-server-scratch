FROM redis:8.2

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update &&\
    apt install -yqq curl dnsutils net-tools iputils-ping mc procps screen telnet tshark wget

COPY entrypoint-custom.sh /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/entrypoint-custom.sh"]
