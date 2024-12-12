FROM golang:1.23 AS BUILDER

LABEL maintainer="Eugene Klimov <bloodjazman@gmail.com>"

COPY . /go/src/github.com/ClickHouse/clickhouse_exporter

WORKDIR /go/src/github.com/ClickHouse/clickhouse_exporter

RUN make init
RUN make all

FROM alpine:latest

COPY --from=BUILDER /go/bin/clickhouse_exporter /usr/local/bin/clickhouse_exporter
RUN apk update && apk add ca-certificates libc6-compat && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/clickhouse_exporter"]
CMD ["-scrape_uri=http://localhost:8123"]
USER   nobody
EXPOSE 9116
