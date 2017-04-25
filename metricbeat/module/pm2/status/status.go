package status

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/helper"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/parse"
)

const (
	// defaultScheme is the default scheme to use when it is not specified in
	// the host config.
	defaultScheme = "http"

	// defaultPath is the default path.
	defaultPath = "/"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   defaultPath,
	}.Build()
)

func init() {
	if err := mb.Registry.AddMetricSet("pm2", "status", New, hostParser); err != nil {
		panic(err)
	}
}

// MetricSet for fetching Nginx stub status.
type MetricSet struct {
	mb.BaseMetricSet
	http           *helper.HTTP
	includeNameMap map[string]bool
	excludeNameMap map[string]bool
}

// New creates new instance of MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	config := struct {
		IncludeNames []string `config:"include_names"`
		ExcludeNames []string `config:"exclude_names"`
	}{}

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	includeNameMap := make(map[string]bool)
	for _, name := range config.IncludeNames {
		includeNameMap[name] = true
	}

	excludeNameMap := make(map[string]bool)
	for _, name := range config.ExcludeNames {
		excludeNameMap[name] = true
	}

	return &MetricSet{
		BaseMetricSet:  base,
		http:           helper.NewHTTP(base),
		includeNameMap: includeNameMap,
		excludeNameMap: excludeNameMap,
	}, nil
}

// Fetch makes an HTTP request to fetch status metrics from the stubstatus endpoint.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	content, err := m.http.FetchContent()
	if err != nil {
		return nil, err
	}

	return eventMapping(content, m)
}
