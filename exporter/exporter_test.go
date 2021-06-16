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

func TestParseNumber(t *testing.T) {
	type testCase struct {
		in  string
		out float64
	}

	testCases := []testCase{
		{in: "1", out: 1},
		{in: "1.1", out: 1.1},
	}

	for _, tc := range testCases {
		out, err := parseNumber(tc.in)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if out != tc.out {
			t.Fatalf("wrong output: %f, expected %f", out, tc.out)
		}
	}
}
