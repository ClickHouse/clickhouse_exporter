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
	TABLE_METRIC_EXPORTER_QUERY = `select database, name as table, engine, total_rows, total_bytes, parts from system.tables {FILTER_CLAUSE}`
)

type TableMetricsExporter struct {
	Namespace string
	QueryURI  string
}

func NewTableMetricsExporter(uri url.URL, namespace string, yamlconfig yaml.YamlConfig) TableMetricsExporter {

	filter_calause := queryparser.ParseYamlConfigToQueryFilter(yamlconfig)
	query := strings.Replace(TABLE_METRIC_EXPORTER_QUERY, "{FILTER_CLAUSE}", filter_calause, 1)
	log.Printf("table exporter query: %v", query)

	url_values := uri.Query()
	metricsURI := uri
	url_values.Set("query", query)
	metricsURI.RawQuery = url_values.Encode()

	return TableMetricsExporter{
		QueryURI:  metricsURI.String(),
		Namespace: namespace,
	}
}

func (e *TableMetricsExporter) Scrap(clickConn clickhouse.ClickhouseConn, ch chan<- prometheus.Metric) error {
	table_metrics, err := e.parseResponse(clickConn)
	if err != nil {
		return fmt.Errorf("error scraping clickhouse url %v: %v", e.QueryURI, err)
	}
	e.collect(table_metrics, ch)

	return nil
}

type TableMetricsResult struct {
	database    string
	table       string
	engine      string
	total_rows  int
	total_bytes int
	parts       int
}

func (e *TableMetricsExporter) parseResponse(clickConn clickhouse.ClickhouseConn) ([]TableMetricsResult, error) {
	data, err := clickConn.ExcecuteQuery(e.QueryURI)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results = make([]TableMetricsResult, 0)

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) != 6 {
			return nil, fmt.Errorf("parseQueryResponse: unexpected %d line: %s", i, line)
		}
		database := strings.TrimSpace(fields[0])
		table := strings.TrimSpace(fields[1])
		engine := strings.TrimSpace(fields[2])

		total_rows, err := strconv.Atoi(strings.TrimSpace(fields[3]))
		if err != nil {
			return nil, err
		}

		total_bytes, err := strconv.Atoi(strings.TrimSpace(fields[4]))
		if err != nil {
			return nil, err
		}

		parts, err := strconv.Atoi(strings.TrimSpace(fields[5]))
		if err != nil {
			return nil, err
		}

		results = append(results, TableMetricsResult{
			database:    database,
			table:       table,
			engine:      engine,
			total_rows:  total_rows,
			total_bytes: total_bytes,
			parts:       parts,
		})
	}

	return results, nil

}

func (e *TableMetricsExporter) collect(resultLines []TableMetricsResult, ch chan<- prometheus.Metric) {

	for _, query_metrics := range resultLines {

		metric_label := []string{"database", "table", "engine"}

		newTotalRows := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_rows",
			Help:      "number of rows of a table",
		}, metric_label).WithLabelValues(query_metrics.database, query_metrics.table, query_metrics.engine)
		newTotalRows.Set(float64(query_metrics.total_rows))
		newTotalRows.Collect(ch)

		newTotalBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_bytes",
			Help:      "table compressed bytes volume",
		}, metric_label).WithLabelValues(query_metrics.database, query_metrics.table, query_metrics.engine)
		newTotalBytes.Set(float64(query_metrics.total_bytes))
		newTotalBytes.Collect(ch)

		newParts := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "table_parts",
			Help:      "number of current table partitions",
		}, metric_label).WithLabelValues(query_metrics.database, query_metrics.table, query_metrics.engine)
		newParts.Set(float64(query_metrics.parts))
		newParts.Collect(ch)
	}

}
