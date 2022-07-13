FROM golang:1.15-buster AS BUILD

LABEL maintainer="Igor Petrenko"

COPY . /go/src/github.com/f1yegor/clickhouse_exporter

WORKDIR /go/src/github.com/f1yegor/clickhouse_exporter

RUN make init && make

FROM frolvlad/alpine-glibc:alpine-3.16

COPY --from=BUILD /go/bin/clickhouse_exporter /usr/local/bin/clickhouse_exporter

RUN apk update && apk add ca-certificates curl bash && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/clickhouse_exporter"]

CMD ["-scrape_uri=http://localhost:8123"]

EXPOSE 9116