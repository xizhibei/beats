// +build integration

package status

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/elastic/beats/metricbeat/module/pm2"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	f := mbtest.NewEventFetcher(t, getConfig())
	event, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)

	// Check event fields
	totalMem := event["monit"].(common.MapStr)["total_mem"].(int64)
	assert.True(t, totalMem >= 0)

	freeMem := event["monit"].(common.MapStr)["free_mem"].(int64)
	assert.True(t, freeMem > 0)
}

func TestData(t *testing.T) {
	f := mbtest.NewEventFetcher(t, getConfig())
	err := mbtest.WriteEvent(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "pm2",
		"metricsets": []string{"status"},
		"hosts":      []string{pm2.GetEnvHost() + ":" + pm2.GetEnvPort()},
	}
}
