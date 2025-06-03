package configs

import (
	"flag"
	"os"
)

type Configuration struct {
	ListeningAddress *string
	MetricsEndpoint  *string
	ClickhouseOnly   *bool
	Insecure         *bool

	ClickhouseScrapeURI string
	User                string
	Password            string

	QueryFiltersPath string
}

func LoadConfigs() Configuration {
	configs := Configuration{
		ListeningAddress: flag.String("telemetry.address", ":9116", "Address on which to expose metrics."),
		MetricsEndpoint:  flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics."),
		ClickhouseOnly:   flag.Bool("clickhouse_only", false, "Expose only Clickhouse metrics, not metrics from the exporter itself"),
		Insecure:         flag.Bool("insecure", true, "Ignore server certificate if using https"),

		ClickhouseScrapeURI: getEnv("CLICKHOUSE_URI", "http://127.0.0.1:8123"),
		User:                getEnv("CLICKHOUSE_USER", "user"),
		Password:            getEnv("CLICKHOUSE_PASSWORD", "pass"),

		QueryFiltersPath: getEnv("QUERY_FILTERS_PATH", "./conf/query-filters.yaml"),
	}

	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()

	return configs
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
