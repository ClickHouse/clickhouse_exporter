# Clickhouse Exporter for Prometheus

This is a simple server that periodically scrapes [ClickHouse](https://clickhouse.com/) stats and exports them via HTTP for [Prometheus](https://prometheus.io/)
consumption.

Exporter could used only for old ClickHouse versions, modern versions have embedded prometheus endpoint.
Look details https://clickhouse.com/docs/en/operations/server-configuration-parameters/settings#server_configuration_parameters-prometheus

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
CLICKHOUSE_URI
CLICKHOUSE_USER
CLICKHOUSE_PASSWORD
```

## Build Docker image
```
docker build . -t clickhouse-exporter \
    --build-arg BUILD_HTTP_PROXY=http://proxy-host:port \
    --build-arg BUILD_HTTPS_PROXY=http://proxy-host:port
```

## Using Docker

```
docker run -d -p 9116:9116 clickhouse-exporter
```
## Sample dashboard
Grafana dashboard could be a start for inspiration https://grafana.com/grafana/dashboards/882-clickhouse

## Using with Prometheus
just add it like this to your prometheus.yaml
```yaml
global:
  scrape_interval:     1s
  evaluation_interval: 1s

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ['127.0.0.1:9090']

  - job_name: clickhouse-exporter
    static_configs:
      - targets: ['127.0.0.1:9116']
```