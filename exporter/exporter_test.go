package exporter

import (
	"net/url"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestScrape(t *testing.T) {
	clickhouseUrl, err := url.Parse("http://127.0.0.1:8123/")
	if err != nil {
		t.Fatal(err)
	}
	exporter := NewExporter(*clickhouseUrl, false, "", "")

	t.Run("Describe", func(t *testing.T) {
		ch := make(chan *prometheus.Desc)
		go func() {
			exporter.Describe(ch)
			close(ch)
		}()

		for range ch {
		}
	})

	t.Run("Collect", func(t *testing.T) {
		ch := make(chan prometheus.Metric)
		var err error
		go func() {
			err = exporter.collect(ch)
			if err != nil {
				panic("failed")
			}
			close(ch)
		}()

		for range ch {
		}
	})
}
