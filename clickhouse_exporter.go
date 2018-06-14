package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"

	"github.com/f1yegor/clickhouse_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

func getCredentials() (bool, string, string) {
	user, userPresent := os.LookupEnv("CLICKHOUSE_USER")
	password, passwordPresent := os.LookupEnv("CLICKHOUSE_PASSWORD")
	return userPresent && passwordPresent, user, password
}

var (
	listeningAddress                   = flag.String("telemetry.address", ":9116", "Address on which to expose metrics.")
	metricsEndpoint                    = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	clickhouseScrapeURI                = flag.String("scrape_uri", "http://localhost:8123/", "URI to clickhouse http endpoint")
	insecure                           = flag.Bool("insecure", true, "Ignore server certificate if using https")
	credentialsPresent, user, password = getCredentials()
)

func main() {
	flag.Parse()

	uri, err := url.Parse(*clickhouseScrapeURI)
	if err != nil {
		log.Fatal(err)
	}
	e := exporter.NewExporter(*uri, *insecure, credentialsPresent, user, password)
	prometheus.MustRegister(e)

	log.Printf("Starting Server: %s", *listeningAddress)
	http.Handle(*metricsEndpoint, prometheus.Handler())
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
