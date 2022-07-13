FROM --platform=$BUILDPLATFORM golang:1.16 AS BUILD

LABEL maintainer="Roman Tkalenko"

COPY . /go/src/github.com/ClickHouse/clickhouse_exporter

WORKDIR /go/src/github.com/ClickHouse/clickhouse_exporter

ARG TARGETARCH

RUN GOOS=linux GOARCH=$TARGETARCH go build -o /go/bin/clickhouse_exporter


FROM pingcap/alpine-glibc:alpine-3.14

COPY --from=BUILD /go/bin/clickhouse_exporter /usr/local/bin/clickhouse_exporter
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/clickhouse_exporter"]
CMD ["-scrape_uri=http://localhost:8123"]
EXPOSE 9116