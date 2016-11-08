package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	clickhouseMetrics = `Query   1
Merge   0
ReplicatedFetch 0
ReplicatedSend  0
ReplicatedChecks        0
BackgroundPoolTask      0
DiskSpaceReservedForMerge       0
DistributedSend 0
QueryPreempted  0
TCPConnection   0
HTTPConnection  1
InterserverConnection   0
OpenFileForRead 0
OpenFileForWrite        0
Read    1
Write   0
SendExternalTables      0
QueryThread     0
ReadonlyReplica 0
MemoryTracking  8704
`
	metricCount = 20
)

func TestClickhouseStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(clickhouseMetrics))
	})
	server := httptest.NewServer(handler)

	e := NewExporter(server.URL)
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()
	// because asks 3 tables
	for i := 1; i <= 3*metricCount; i++ {
		m := <-ch
		if m == nil {
			t.Error("expected metric but got nil")
		}
	}
	if <-ch != nil {
		t.Error("expected closed channel")
	}
}
