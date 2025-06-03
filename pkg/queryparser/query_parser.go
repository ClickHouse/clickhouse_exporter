package queryparser

import (
	"strings"

	"github.com/ClickHouse/clickhouse_exporter/pkg/yaml"
)

func ParseYamlConfigToQueryFilter(yamlObject yaml.YamlConfig) string {
	containsFilters := yamlObject.Contains("filters")
	if !containsFilters {
		return ""
	}
	query_filter_clause := makeFilterCaluse(yamlObject)
	return query_filter_clause
}

func makeFilterCaluse(yamlObject yaml.YamlConfig) string {

	filterValue := yamlObject.GetValue("filters")

	switch val := filterValue.(type) {
	case []interface{}:
		{
			return makeFilterFromInterfaces(val)

		}
	case string:
		{
			return makeFilterFromString(val)
		}
	case nil:
		{
			return ""
		}
	default:
		panic("error: this is not the proper format for setting filters. use list or single string")
	}
}

func makeFilterFromInterfaces(interfaces []interface{}) string {
	var filtersString []string
	for _, interface_unit := range interfaces {
		filtersString = append(filtersString, interface_unit.(string))
	}

	return "WHERE\n" + strings.Join(filtersString, " AND\n")
}

func makeFilterFromString(s string) string {
	return "WHERE\n" + s
}
