FROM ubuntu:24.04

ARG TARGETPLATFORM

RUN useradd -m privatebin && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY $TARGETPLATFORM/privatebin /usr/local/bin/privatebin
RUN chmod +x /usr/local/bin/privatebin

USER privatebin

ENTRYPOINT ["/usr/local/bin/privatebin"]
