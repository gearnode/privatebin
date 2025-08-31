FROM ubuntu:22.04

LABEL org.opencontainers.image.source="https://github.com/gearnode/privatebin"
LABEL org.opencontainers.image.licenses="ISC"

RUN useradd -m privatebin && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY privatebin /usr/local/bin/privatebin
RUN chmod +x /usr/local/bin/privatebin

USER privatebin

ENTRYPOINT ["/usr/local/bin/privatebin"]