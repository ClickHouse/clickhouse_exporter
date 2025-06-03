FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

ARG BUILD_HTTP_PROXY
ARG BUILD_HTTPS_PROXY
ENV HTTP_PROXY=${BUILD_HTTP_PROXY}
ENV HTTPS_PROXY=${BUILD_HTTPS_PROXY}

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o clickhouse_exporter ./cmd/clickhouse_exporter


FROM alpine:3.21.3

ARG BUILD_HTTP_PROXY
ARG BUILD_HTTPS_PROXY
ENV HTTP_PROXY=${BUILD_HTTP_PROXY}
ENV HTTPS_PROXY=${BUILD_HTTPS_PROXY}

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/clickhouse_exporter /opt/clickhouse_exporter/runner

COPY ./conf /opt/clickhouse_exporter/conf
ENV QUERY_FILTERS_PATH=/opt/clickhouse_exporter/conf/query-filters.yaml

USER nobody

ARG PORT
EXPOSE ${PORT}

ENTRYPOINT ["/opt/clickhouse_exporter/runner"]