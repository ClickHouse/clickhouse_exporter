package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"

	"github.com/f1yegor/clickhouse_exporter/exporter"
	"github.com/peterbourgon/ff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
)

var (
	fs                  = flag.NewFlagSet("clickhouse-exporter", flag.ExitOnError)
	listeningAddress    = fs.String("telemetry.address", ":9116", "Address on which to expose metrics. Override environment (CLICKHOUSE_TELEMETRY_ADDRESS)")
	metricsEndpoint     = fs.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics. Override environment (CLICKHOUSE_TELEMETRY_ENDPOINT)")
	clickhouseScrapeURI = fs.String("scrape_uri", "http://localhost:8123/", "URI to clickhouse http endpoint. Override environment (CLICKHOUSE_SCRAPE_URI)")
	insecure            = fs.Bool("insecure", true, "Ignore server certificate if using https. Override environment (CLICKHOUSE_INSECURE)")
	user                = fs.String("user", "", "Clickhouse user. Override environment (CLICKHOUSE_USER)")
	password            = fs.String("password", "", "Clickhouse password. Override environment (CLICKHOUSE_PASSWORD)")
)

func main() {
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("CLICKHOUSE")); err != nil {
		panic(err)
	}

	uri, err := url.Parse(*clickhouseScrapeURI)
	if err != nil {
		log.Fatal(err)
	}
	e := exporter.NewExporter(*uri, *insecure, *user, *password)
	prometheus.MustRegister(e)

	log.Printf("Starting Server: %s", *listeningAddress)
	http.Handle(*metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Clickhouse Exporter</title></head>
			<body>
			<h1>Clickhouse Exporter</h1>
			<p><a href="` + *metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
