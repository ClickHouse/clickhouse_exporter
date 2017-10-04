# Clickhouse Exporter for Prometheus

[![Build Status](https://travis-ci.org/f1yegor/clickhouse_exporter.svg?branch=master)](https://travis-ci.org/f1yegor/clickhouse_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/f1yegor/clickhouse_exporter)](https://goreportcard.com/report/github.com/f1yegor/clickhouse_exporter)

This is a simple server that periodically scrapes ClickHouse(https://clickhouse.yandex/) stats and exports them via HTTP for Prometheus(https://prometheus.io/)
consumption.

To run it:

```bash
./clickhouse_exporter [flags]
```

Help on flags:
```bash
./clickhouse_exporter --help
```

Credentials(if not default):

via environment variables
```
CLICKHOUSE_USER
CLICKHOUSE_PASSWORD
```

## Using Docker

```
docker run -d -p 9116:9116 f1yegor/clickhouse-exporter -scrape_uri=http://clickhouse.service.consul:8123/
```
## Sample dashboard
Grafana dashboard could be a start for inspiration https://grafana.net/dashboards/882
