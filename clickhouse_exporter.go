package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/f1yegor/clickhouse_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listeningAddress    = flag.String("telemetry.address", ":9116", "Address on which to expose metrics.")
	metricsEndpoint     = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	readTimeoutSeconds  = flag.Int("telemetry.read_timeout_seconds", 5, "Maximum duration before timing out read of the request, and closing idle connections.")
	writeTimeoutSeconds = flag.Int("telemetry.write_timeout_seconds", 10, "Maximum duration before timing out write of the response.")
	clickhouseScrapeURI = flag.String("scrape_uri", "http://localhost:8123/", "URI to clickhouse http endpoint")
	insecure            = flag.Bool("insecure", true, "Ignore server certificate if using https")
	user                = os.Getenv("CLICKHOUSE_USER")
	password            = os.Getenv("CLICKHOUSE_PASSWORD")
)

func main() {
	flag.Parse()

	uri, err := url.Parse(*clickhouseScrapeURI)
	if err != nil {
		log.Fatal(err)
	}
	e := exporter.NewExporter(*uri, *insecure, user, password)
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

	srv := &http.Server{
		Addr:         *listeningAddress,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(*readTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(*writeTimeoutSeconds) * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
