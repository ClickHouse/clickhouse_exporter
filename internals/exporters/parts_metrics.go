package exporters

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ClickHouse/clickhouse_exporter/pkg/clickhouse"
	"github.com/ClickHouse/clickhouse_exporter/pkg/queryparser"
	"github.com/ClickHouse/clickhouse_exporter/pkg/yaml"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

const (
	PARTS_METRIC_EXPORTER_QUERY = `select database, table, sum(bytes) as bytes, count() as parts, sum(rows) as rows from system.parts {FILTER_CLAUSE} group by database, table`
)

type PartsMetricsExporter struct {
	Namespace string
	QueryURI  string
}

func NewPartsMetricsExporter(uri url.URL, namespace string, yamlconfig yaml.YamlConfig) PartsMetricsExporter {

	filter_calause := queryparser.ParseYamlConfigToQueryFilter(yamlconfig)
	query := strings.Replace(PARTS_METRIC_EXPORTER_QUERY, "{FILTER_CLAUSE}", filter_calause, 1)
	log.Printf("parts exporter query: %v", query)

	url_values := uri.Query()
	metricsURI := uri
	url_values.Set("query", query)
	metricsURI.RawQuery = url_values.Encode()

	return PartsMetricsExporter{
		QueryURI:  metricsURI.String(),
		Namespace: namespace,
	}
}

func (e *PartsMetricsExporter) Scrap(clickConn clickhouse.ClickhouseConn, ch chan<- prometheus.Metric) error {
	parts, err := e.parseResponse(clickConn)
	if err != nil {
		return fmt.Errorf("error scraping clickhouse url %v: %v", e.QueryURI, err)
	}
	e.collect(parts, ch)
	return nil
}

type PartsResult struct {
	database string
	table    string
	bytes    int
	parts    int
	rows     int
}

func (e *PartsMetricsExporter) parseResponse(clickConn clickhouse.ClickhouseConn) ([]PartsResult, error) {
	data, err := clickConn.ExcecuteQuery(e.QueryURI)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results = make([]PartsResult, 0)

	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) != 5 {
			return nil, fmt.Errorf("parsePartsResponse: unexpected %d line: %s", i, line)
		}
		database := strings.TrimSpace(parts[0])
		table := strings.TrimSpace(parts[1])

		bytes, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil {
			return nil, err
		}

		count, err := strconv.Atoi(strings.TrimSpace(parts[3]))
		if err != nil {
			return nil, err
		}

		rows, err := strconv.Atoi(strings.TrimSpace(parts[4]))
		if err != nil {
			return nil, err
		}

		results = append(results, PartsResult{database, table, bytes, count, rows})
	}

	return results, nil
}

func (e *PartsMetricsExporter) collect(resultLines []PartsResult, ch chan<- prometheus.Metric) {
	for _, part := range resultLines {
		newBytesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_parts_bytes",
			Help:      "Table size in bytes",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newBytesMetric.Set(float64(part.bytes))
		newBytesMetric.Collect(ch)

		newCountMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_parts_count",
			Help:      "Number of parts of the table",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newCountMetric.Set(float64(part.parts))
		newCountMetric.Collect(ch)

		newRowsMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_parts_rows",
			Help:      "Number of rows in the table",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newRowsMetric.Set(float64(part.rows))
		newRowsMetric.Collect(ch)
	}

}
