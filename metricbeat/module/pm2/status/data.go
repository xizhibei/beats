package status

import (
	"encoding/json"

	"github.com/elastic/beats/libbeat/common"
)

type PM2Env struct {
	ExecMode         string `json:"exec_mode"`
	PMUptime         int64  `json:"pm_uptime"`
	Status           string `json:"status"`
	NodeVersion      string `json:"node_version"`
	Script           string `json:"script"`
	Restarts         int32  `json:"restart_time"`
	UnstableRestarts int32  `json:"unstable_restarts"`
}

type SystemInfo struct {
	HostName string `json:"hostname"`
	Uptime   int32  `json:"uptime"`
}

type Monit struct {
	Memory int32 `json:"memory"`
	CPU    int32 `json:"cpu"`
}

type Process struct {
	PId    int32  `json:"pid"`
	Name   string `json:"name"`
	PMId   int32  `json:"pm_id"`
	Monit  Monit  `json:"monit"`
	PM2Env PM2Env `json:"pm2_env"`
}

type Output struct {
	Processes  []Process  `json:"processes"`
	SystemInfo SystemInfo `json:"system_info"`
}

// Map body to MapStr
func eventMapping(content []byte, m *MetricSet) ([]common.MapStr, error) {
	var o Output

	err := json.Unmarshal(content, &o)
	if err != nil {
		return nil, err
	}

	events := []common.MapStr{}

	for _, Process := range o.Processes {
		if len(m.includeNameMap) > 0 && !m.includeNameMap[Process.Name] {
			continue
		}
		if m.excludeNameMap[Process.Name] {
			continue
		}
		event := common.MapStr{
			"pid":               Process.PId,
			"name":              Process.Name,
			"pm_id":             Process.PMId,
			"memory":            Process.Monit.Memory,
			"cpu":               Process.Monit.CPU,
			"status":            Process.PM2Env.Status,
			"pm_uptime":         Process.PM2Env.PMUptime,
			"restarts":          Process.PM2Env.Restarts,
			"unstable_restarts": Process.PM2Env.UnstableRestarts,
			"node_version":      Process.PM2Env.NodeVersion,
			"exec_mode":         Process.PM2Env.ExecMode,
			"hostname":          o.SystemInfo.HostName,
		}

		events = append(events, event)
	}

	return events, nil
}
