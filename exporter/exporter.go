package exporter

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	namespace = "clickhouse" // For Prometheus metrics.
)

// Exporter collects clickhouse stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	metricsURI      string
	asyncMetricsURI string
	eventsURI       string
	partsURI        string
	mutex           sync.RWMutex
	client          *http.Client

	scrapeFailures prometheus.Counter

	gauges   []*prometheus.GaugeVec
	counters []*prometheus.CounterVec

	user     string
	password string
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri url.URL, insecure bool, user, password string) *Exporter {
	q := uri.Query()
	metricsURI := uri
	q.Set("query", "select metric, value from system.metrics")
	metricsURI.RawQuery = q.Encode()

	asyncMetricsURI := uri
	q.Set("query", "select * from system.asynchronous_metrics")
	asyncMetricsURI.RawQuery = q.Encode()

	eventsURI := uri
	q.Set("query", "select event, value from system.events")
	eventsURI.RawQuery = q.Encode()

	partsURI := uri
	q.Set("query", "select database, table, sum(bytes) as bytes, count() as parts, sum(rows) as rows from system.parts where active = 1 group by database, table")
	partsURI.RawQuery = q.Encode()
	
	return &Exporter{
		metricsURI:      metricsURI.String(),
		asyncMetricsURI: asyncMetricsURI.String(),
		eventsURI:       eventsURI.String(),
		partsURI:        partsURI.String(),
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrape_failures_total",
			Help:      "Number of errors while scraping clickhouse.",
		}),
		gauges:   make([]*prometheus.GaugeVec, 0, 20),
		counters: make([]*prometheus.CounterVec, 0, 20),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		},
		user:     user,
		password: password,
	}
}

// Describe describes all the metrics ever exported by the clickhouse exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
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

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	metrics, err := e.parseKeyValueResponse(e.metricsURI)
	if err != nil {
		return fmt.Errorf("Error scraping clickhouse url %v: %v", e.metricsURI, err)
	}

	for _, m := range metrics {
		newMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      metricName(m.key),
			Help:      "Number of " + m.key + " currently processed",
		}, []string{}).WithLabelValues()
		newMetric.Set(float64(m.value))
		newMetric.Collect(ch)
	}

	asyncMetrics, err := e.parseKeyValueResponse(e.asyncMetricsURI)
	if err != nil {
		return fmt.Errorf("Error scraping clickhouse url %v: %v", e.asyncMetricsURI, err)
	}

	for _, am := range asyncMetrics {
		newMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      metricName(am.key),
			Help:      "Number of " + am.key + " async processed",
		}, []string{}).WithLabelValues()
		newMetric.Set(float64(am.value))
		newMetric.Collect(ch)
	}

	events, err := e.parseKeyValueResponse(e.eventsURI)
	if err != nil {
		return fmt.Errorf("Error scraping clickhouse url %v: %v", e.eventsURI, err)
	}

	for _, ev := range events {
		newMetric, _ := prometheus.NewConstMetric(
			prometheus.NewDesc(
				namespace+"_"+metricName(ev.key)+"_total",
				"Number of "+ev.key+" total processed", []string{}, nil),
			prometheus.CounterValue, float64(ev.value))
		ch <- newMetric
	}

	parts, err := e.parsePartsResponse(e.partsURI)
	if err != nil {
		return fmt.Errorf("Error scraping clickhouse url %v: %v", e.partsURI, err)
	}

	for _, part := range parts {
		newBytesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "table_parts_bytes",
			Help:      "Table size in bytes",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newBytesMetric.Set(float64(part.bytes))
		newBytesMetric.Collect(ch)

		newCountMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "table_parts_count",
			Help:      "Number of parts of the table",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newCountMetric.Set(float64(part.parts))
		newCountMetric.Collect(ch)

		newRowsMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "table_parts_rows",
			Help:      "Number of rows in the table",
		}, []string{"database", "table"}).WithLabelValues(part.database, part.table)
		newRowsMetric.Set(float64(part.rows))
		newRowsMetric.Collect(ch)
	}

	return nil
}

func (e *Exporter) handleResponse(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	if e.user != "" && e.password != "" {
		req.Header.Set("X-ClickHouse-User", e.user)
		req.Header.Set("X-ClickHouse-Key", e.password)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error scraping clickhouse: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return nil, fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}
	
	if strings.Contains(uri, "system.metrics") || strings.Contains(uri, "system.events") {
		data = DataFix(data)
	}
	return data, nil
}

func DataFix(data []byte) []byte {
	datastr := string(data)
	var fixedData string
	for _, line := range strings.Split(datastr, "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fixedData = fixedData + fmt.Sprintf("%v %v\n", parts[0], parts[1])
		}
	}
	return []byte(fixedData)
}

type lineResult struct {
	key   string
	value int
}

func (e *Exporter) parseKeyValueResponse(uri string) ([]lineResult, error) {
	data, err := e.handleResponse(uri)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results []lineResult = make([]lineResult, 0)

	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) != 2 {
			return nil, fmt.Errorf("parseKeyValueResponse: unexpected %d line: %s", i, line)
		}
		k := strings.TrimSpace(parts[0])
		v, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}
		results = append(results, lineResult{k, v})

	}
	return results, nil
}

type partsResult struct {
	database string
	table    string
	bytes    int
	parts    int
	rows     int
}

func (e *Exporter) parsePartsResponse(uri string) ([]partsResult, error) {
	data, err := e.handleResponse(uri)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results []partsResult = make([]partsResult, 0)

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

		results = append(results, partsResult{database, table, bytes, count, rows})
	}

	return results, nil
}

// Collect fetches the stats from configured clickhouse location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Printf("Error scraping clickhouse: %s", err)
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
	}
	// Reset metrics.
	for _, vec := range e.gauges {
		vec.Reset()
	}

	for _, vec := range e.counters {
		vec.Reset()
	}

	for _, vec := range e.gauges {
		vec.Collect(ch)
	}

	for _, vec := range e.counters {
		vec.Collect(ch)
	}

	return
}

func metricName(in string) string {
	out := toSnake(in)
	return strings.Replace(out, ".", "_", -1)
}

// toSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// check interface
var _ prometheus.Collector = (*Exporter)(nil)
