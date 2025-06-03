package exporters

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ClickHouse/clickhouse_exporter/internals/util"
	"github.com/ClickHouse/clickhouse_exporter/pkg/clickhouse"
	"github.com/ClickHouse/clickhouse_exporter/pkg/queryparser"
	"github.com/ClickHouse/clickhouse_exporter/pkg/yaml"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

const (
	ASYNC_METRIC_EXPORTER_QUERY = `
	select replaceRegexpAll(toString(metric), '-', '_') AS metric, value from system.asynchronous_metrics {FILTER_CLAUSE}`
)

type AsyncMetricsExporter struct {
	Namespace string
	QueryURI  string
}

func NewAsyncMetricsExporter(uri url.URL, namespace string, yamlconfig yaml.YamlConfig) AsyncMetricsExporter {

	filter_calause := queryparser.ParseYamlConfigToQueryFilter(yamlconfig)
	query := strings.Replace(ASYNC_METRIC_EXPORTER_QUERY, "{FILTER_CLAUSE}", filter_calause, 1)
	log.Printf("async exporter query: %v", query)

	url_values := uri.Query()
	metricsURI := uri
	url_values.Set("query", query)
	metricsURI.RawQuery = url_values.Encode()

	return AsyncMetricsExporter{
		QueryURI:  metricsURI.String(),
		Namespace: namespace,
	}
}

func (e *AsyncMetricsExporter) Scrap(clickConn clickhouse.ClickhouseConn, ch chan<- prometheus.Metric) error {
	asyncMetrics, err := util.ParseKeyValueResponse(e.QueryURI, clickConn)
	if err != nil {
		return fmt.Errorf("error scraping clickhouse url %v: %v", e.QueryURI, err)
	}
	e.collect(asyncMetrics, ch)
	return nil
}

func (e *AsyncMetricsExporter) collect(resultLines []util.LineResult, ch chan<- prometheus.Metric) {
	for _, am := range resultLines {
		newMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      util.GetMetricName(am.Key),
			Help:      "Number of " + am.Key + " async processed",
		}, []string{}).WithLabelValues()
		newMetric.Set(am.Value)
		newMetric.Collect(ch)
	}
}
