# Clickhouse Exporter for Prometheus

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

## Using Docker

```
docker run -d -p 9116:9116 f1yegor/clickhouse-exporter -scrape_uri=http://clickhouse.service.consul:8123/
```
## Sample dashboard
Grafana dashboard could be a start for inspiration https://grafana.net/dashboards/882
