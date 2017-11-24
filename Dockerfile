FROM alpine

MAINTAINER Robert Buck <bob@continuul.io>

ENV GOSU_VERSION 1.10

RUN apk update && apk upgrade && \
    apk add --no-cache curl git && \
    curl -o /usr/local/bin/gosu -sSL "https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-amd64" && \
    chmod +x /usr/local/bin/gosu && \
    apk del curl && \
    rm -rf /var/cache/apk/*

RUN addgroup continuul && \
    adduser -S -G continuul continuul

RUN mkdir -p /var/lib/continuul && \
    mkdir -p /var/log/continuul && \
    mkdir -p /etc/continuul && \
    chown -R continuul:continuul /var/lib/continuul /var/log/continuul

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

COPY on /usr/local/bin

VOLUME /var/lib/continuul

EXPOSE 9700

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["agent", "--bind", "0.0.0.0", "--snapshot", "/var/tmp/serf.snapshot"]
