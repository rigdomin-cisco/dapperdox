FROM alpine:3.10

RUN apk add --no-cache \
        ca-certificates \
        bash \
    && rm -f /var/cache/apk/*

ARG VERSION
ENV DAPPERDOX ${VERSION}

COPY bin/dapperdox /usr/local/bin/dapperdox

CMD ["/usr/local/bin/dapperdox"]
