# Clickhouse Exporter for Prometheus

[![Build Status](https://travis-ci.org/f1yegor/clickhouse_exporter.svg?branch=master)](https://travis-ci.org/f1yegor/clickhouse_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/f1yegor/clickhouse_exporter)](https://goreportcard.com/report/github.com/f1yegor/clickhouse_exporter)

This is a simple server that periodically scrapes ClickHouse(https://clickhouse.yandex/) stats and exports them via HTTP for Prometheus(https://prometheus.io/)
consumption.

To run it:

```bash
./clickhouse_exporter [flags]
```

Flags are also configurable using environment variables. See the usage with:

```bash
./clickhouse_exporter --help
```

```
Usage of clickhouse-exporter:
  -insecure
    	Ignore server certificate if using https. Override environment (CLICKHOUSE_INSECURE) (default true)
  -password string
    	Clickhouse password. Override environment (CLICKHOUSE_PASSWORD)
  -scrape_uri string
    	URI to clickhouse http endpoint. Override environment (CLICKHOUSE_SCRAPE_URI) (default "http://localhost:8123/")
  -telemetry.address string
    	Address on which to expose metrics. Override environment (CLICKHOUSE_TELEMETRY_ADDRESS) (default ":9116")
  -telemetry.endpoint string
    	Path under which to expose metrics. Override environment (CLICKHOUSE_TELEMETRY_ENDPOINT) (default "/metrics")
  -user string
    	Clickhouse user. Override environment (CLICKHOUSE_USER)
```

## Using Docker

```
docker run -d -p 9116:9116 f1yegor/clickhouse-exporter -scrape_uri=http://clickhouse.service.consul:8123/
```

## Sample dashboard

Grafana dashboard could be a start for inspiration https://grafana.net/dashboards/882
