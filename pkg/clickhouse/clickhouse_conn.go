package clickhouse

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ClickhouseConn struct {
	Client   *http.Client
	User     string
	Password string
}

func (e *ClickhouseConn) ExcecuteQuery(query string) ([]byte, error) {
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	if e.User != "" && e.Password != "" {
		req.Header.Set("X-ClickHouse-User", e.User)
		req.Header.Set("X-ClickHouse-Key", e.Password)
	}
	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error scraping clickhouse: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("can't close resp.Body")
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return nil, fmt.Errorf("status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	return data, nil
}
