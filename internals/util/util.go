package util

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/ClickHouse/clickhouse_exporter/pkg/clickhouse"
)

type LineResult struct {
	Key   string
	Value float64
}

func GetMetricName(Key string) string {
	out := toSnake(Key)
	return strings.Replace(out, ".", "_", -1)
}

func ParseNumber(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func ParseKeyValueResponse(uri string, clickConn clickhouse.ClickhouseConn) ([]LineResult, error) {
	data, err := clickConn.ExcecuteQuery(uri)
	if err != nil {
		return nil, err
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	var results = make([]LineResult, 0)

	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) != 2 {
			return nil, fmt.Errorf("parseKeyValueResponse: unexpected %d line: %s", i, line)
		}
		k := strings.TrimSpace(parts[0])
		v, err := ParseNumber(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}
		results = append(results, LineResult{k, v})

	}
	return results, nil
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
