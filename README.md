# Clickhouse Exporter for Prometheus

This is a simple server that periodically scrapes [ClickHouse](https://clickhouse.tech/) stats and exports them via HTTP for [Prometheus](https://prometheus.io/)
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

## Build Docker image
```
docker build . -t clickhouse-exporter
```

## Using Docker

```
docker run -d -p 9116:9116 clickhouse-exporter -scrape_uri=http://clickhouse-url:8123/
```
## Sample dashboard
Grafana dashboard could be a start for inspiration https://grafana.net/dashboards/882
