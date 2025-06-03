package exporter

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/ClickHouse/clickhouse_exporter/internals/exporters"
	"github.com/ClickHouse/clickhouse_exporter/pkg/clickhouse"
	"github.com/ClickHouse/clickhouse_exporter/pkg/configs"
	"github.com/ClickHouse/clickhouse_exporter/pkg/yaml"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

const (
	NAMESPACE = "clickhouse" // For Prometheus metrics.
)

// Exporter collects clickhouse stats from the given URI and exports them using
// the prometheus metrics package.
type ExporterHolder struct {
	basicMetricsExporter exporters.BasicMetricsExporter
	asyncMetricsExporter exporters.AsyncMetricsExporter
	eventMetricsExporter exporters.EventMetricsExporter
	partMetricsExporter  exporters.PartsMetricsExporter
	diskMetricsExporter  exporters.DiskMetricsExporter
	queryMetricsExporter exporters.QueryMetricsExporter
	tableMetricsExporter exporters.TableMetricsExporter

	scrapeFailures prometheus.Counter
	clickConn      clickhouse.ClickhouseConn
}

// NewExporter returns an initialized Exporter.
func NewExporterHolder(configs configs.Configuration) *ExporterHolder {

	uri, err := url.Parse(configs.ClickhouseScrapeURI)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Printf("Scraping %s", configs.ClickhouseScrapeURI)

	queryFilters := yaml.ReadYaml(configs.QueryFiltersPath)

	basicMetricsExporter := exporters.NewBasicMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("basic_exporter"),
	)

	asyncMetricsExporter := exporters.NewAsyncMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("async_exporter"),
	)

	eventMetricsExporter := exporters.NewEventMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("event_exporter"),
	)

	partMetricsExporter := exporters.NewPartsMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("parts_exporter"),
	)

	diskMetricsExporter := exporters.NewDiskMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("disk_exporter"),
	)

	queryMetricsExporter := exporters.NewQueryMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("query_exporter"),
	)

	tableMetricsExporter := exporters.NewTableMetricsExporter(
		*uri,
		NAMESPACE,
		queryFilters.GetMapObject("table_exporter"),
	)

	return &ExporterHolder{
		basicMetricsExporter: basicMetricsExporter,
		asyncMetricsExporter: asyncMetricsExporter,
		eventMetricsExporter: eventMetricsExporter,
		partMetricsExporter:  partMetricsExporter,
		diskMetricsExporter:  diskMetricsExporter,
		queryMetricsExporter: queryMetricsExporter,
		tableMetricsExporter: tableMetricsExporter,
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: NAMESPACE,
			Name:      "exporter_scrape_failures_total",
			Help:      "Number of errors while scraping clickhouse.",
		}),
		clickConn: clickhouse.ClickhouseConn{
			Client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: *configs.Insecure},
				},
				Timeout: 30 * time.Second,
			},
			User:     configs.User,
			Password: configs.Password,
		},
	}
}

// Describe describes all the metrics ever exported by the clickhouse exporter. It
// implements prometheus.Collector.
func (e *ExporterHolder) Describe(ch chan<- *prometheus.Desc) {
	// We cannot know in advance what metrics the exporter will generate
	// from clickhouse. So we use the poor man's describe method: Run a collect
	// and send the descriptors of all the collected metrics.

	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

func (e *ExporterHolder) collect(ch chan<- prometheus.Metric) error {

	e.asyncMetricsExporter.Scrap(e.clickConn, ch)
	e.basicMetricsExporter.Scrap(e.clickConn, ch)
	e.diskMetricsExporter.Scrap(e.clickConn, ch)
	e.eventMetricsExporter.Scrap(e.clickConn, ch)
	e.partMetricsExporter.Scrap(e.clickConn, ch)
	e.queryMetricsExporter.Scrap(e.clickConn, ch)
	e.tableMetricsExporter.Scrap(e.clickConn, ch)

	return nil
}

// Collect fetches the stats from configured clickhouse location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *ExporterHolder) Collect(ch chan<- prometheus.Metric) {
	upValue := 1

	if err := e.collect(ch); err != nil {
		log.Error().Msgf("Error scraping clickhouse: %s", err)
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)

		upValue = 0
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(NAMESPACE, "", "up"),
			"Was the last query of ClickHouse successful.",
			nil, nil,
		),
		prometheus.GaugeValue, float64(upValue),
	)

}

// check interface
var _ prometheus.Collector = (*ExporterHolder)(nil)
