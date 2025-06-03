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
	QUERY_METRIC_EXPORTER_QUERY = `
	SELECT 
		user, type as status, query_kind, arrayJoin(tables) AS table,
		sum(memory_usage) as memory_usage,
		count(*) AS query_num,
		sum(query_duration_ms) as query_duration_ms,
		sum(read_bytes) as read_bytes,
		sum(read_rows) as read_rows,
		sum(written_bytes) as written_bytes,
		sum(written_rows) as written_rows,
		sum(result_bytes) as result_bytes,
		sum(result_rows) as result_rows,
		sum(peak_threads_usage) as peak_threads_usage
	FROM system.query_log 
	{FILTER_CLAUSE}
	GROUP BY user, table, type,query_kind`
)

type QueryMetricsExporter struct {
	Namespace string
	QueryURI  string
}

func NewQueryMetricsExporter(uri url.URL, namespace string, yamlconfig yaml.YamlConfig) QueryMetricsExporter {

	filter_calause := queryparser.ParseYamlConfigToQueryFilter(yamlconfig)
	query := strings.Replace(QUERY_METRIC_EXPORTER_QUERY, "{FILTER_CLAUSE}", filter_calause, 1)
	log.Printf("query exporter query: %v", query)

	url_values := uri.Query()
	metricsURI := uri
	url_values.Set("query", query)
	metricsURI.RawQuery = url_values.Encode()

	return QueryMetricsExporter{
		QueryURI:  metricsURI.String(),
		Namespace: namespace,
	}
}

func (e *QueryMetricsExporter) Scrap(clickConn clickhouse.ClickhouseConn, ch chan<- prometheus.Metric) error {
	query_metrics, err := e.parseResponse(clickConn)
	if err != nil {
		return fmt.Errorf("error scraping clickhouse url %v: %v", e.QueryURI, err)
	}
	e.collect(query_metrics, ch)

	return nil
}

type QueryMetricsResult struct {
	user               string
	query_type         string
	query_kind         string
	table              string
	memory_usage       int
	query_num          int
	query_duration_ms  int
	read_bytes         int
	read_rows          int
	written_bytes      int
	written_rows       int
	result_bytes       int
	result_rows        int
	peak_threads_usage int
}

func (e *QueryMetricsExporter) parseResponse(clickConn clickhouse.ClickhouseConn) ([]QueryMetricsResult, error) {
	data, err := clickConn.ExcecuteQuery(e.QueryURI)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results = make([]QueryMetricsResult, 0)

	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) != 14 {
			return nil, fmt.Errorf("parseQueryResponse: unexpected %d line: %s", i, line)
		}
		user := strings.TrimSpace(parts[0])
		query_type := strings.TrimSpace(parts[1])
		query_kind := strings.TrimSpace(parts[2])
		table := strings.TrimSpace(parts[3])

		memory_usage, err := strconv.Atoi(strings.TrimSpace(parts[4]))
		if err != nil {
			return nil, err
		}

		query_num, err := strconv.Atoi(strings.TrimSpace(parts[5]))
		if err != nil {
			return nil, err
		}

		query_duration_ms, err := strconv.Atoi(strings.TrimSpace(parts[6]))
		if err != nil {
			return nil, err
		}

		read_bytes, err := strconv.Atoi(strings.TrimSpace(parts[7]))
		if err != nil {
			return nil, err
		}

		read_rows, err := strconv.Atoi(strings.TrimSpace(parts[8]))
		if err != nil {
			return nil, err
		}

		written_bytes, err := strconv.Atoi(strings.TrimSpace(parts[9]))
		if err != nil {
			return nil, err
		}

		written_rows, err := strconv.Atoi(strings.TrimSpace(parts[10]))
		if err != nil {
			return nil, err
		}

		result_bytes, err := strconv.Atoi(strings.TrimSpace(parts[11]))
		if err != nil {
			return nil, err
		}

		result_rows, err := strconv.Atoi(strings.TrimSpace(parts[12]))
		if err != nil {
			return nil, err
		}

		peak_threads_usage, err := strconv.Atoi(strings.TrimSpace(parts[13]))
		if err != nil {
			return nil, err
		}

		results = append(results, QueryMetricsResult{
			user:               user,
			query_type:         query_type,
			query_kind:         query_kind,
			table:              table,
			memory_usage:       memory_usage,
			query_num:          query_num,
			query_duration_ms:  query_duration_ms,
			read_bytes:         read_bytes,
			read_rows:          read_rows,
			written_bytes:      written_bytes,
			written_rows:       written_rows,
			result_bytes:       result_bytes,
			result_rows:        result_rows,
			peak_threads_usage: peak_threads_usage,
		})
	}

	return results, nil

}

func (e *QueryMetricsExporter) collect(resultLines []QueryMetricsResult, ch chan<- prometheus.Metric) {

	for _, query_metrics := range resultLines {

		metric_label := []string{"user", "table", "type", "kind"}

		newMemoryUsageMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_memory_usage",
			Help:      "user memory use in bytes",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newMemoryUsageMetric.Set(float64(query_metrics.memory_usage))
		newMemoryUsageMetric.Collect(ch)

		newQueryNumMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_query_num",
			Help:      "Number of Queries that user run",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newQueryNumMetric.Set(float64(query_metrics.query_num))
		newQueryNumMetric.Collect(ch)

		newQueryDurationMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_query_duration_ms",
			Help:      "Duration of Queries in mili seconds",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newQueryDurationMetric.Set(float64(query_metrics.query_duration_ms))
		newQueryDurationMetric.Collect(ch)

		newReadBytesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_read_bytes",
			Help:      "Volume of red rows in bytes",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newReadBytesMetric.Set(float64(query_metrics.read_bytes))
		newReadBytesMetric.Collect(ch)

		newWrittenBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_written_bytes",
			Help:      "Number of bytes that user write",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newWrittenBytes.Set(float64(query_metrics.written_bytes))
		newWrittenBytes.Collect(ch)

		newWrittenRows := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_written_rows",
			Help:      "Number of rows that user write",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newWrittenRows.Set(float64(query_metrics.written_rows))
		newWrittenRows.Collect(ch)

		newResultBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_result_bytes",
			Help:      "Number of result bytes",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newResultBytes.Set(float64(query_metrics.result_bytes))
		newResultBytes.Collect(ch)

		newResultRows := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_result_rows",
			Help:      "Number of result rows",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newResultRows.Set(float64(query_metrics.result_rows))
		newResultRows.Collect(ch)

		newPeakThreadUsage := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.Namespace,
			Name:      "user_peak_thread_usage",
			Help:      "number of threads in the peak",
		}, metric_label).WithLabelValues(query_metrics.user, query_metrics.table, query_metrics.query_type, query_metrics.query_kind)
		newPeakThreadUsage.Set(float64(query_metrics.peak_threads_usage))
		newPeakThreadUsage.Collect(ch)
	}

}
