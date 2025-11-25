package docker

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types/container"
)

func (m *Monitor) watchStats(c *Container) {
	ctx := context.Background()
	stats, err := m.cli.ContainerStats(ctx, c.ID, true)
	if err != nil {
		return
	}

	defer stats.Body.Close()

	dec := json.NewDecoder(stats.Body)
	for {
		var v container.StatsResponse
		if err := dec.Decode(&v); err != nil {
			return
		}

		m.mu.Lock()
		c.Stats = v
		m.mu.Unlock()
		m.events <- Event{EventType: ContainerStatsEventType, ActorID: c.ID}

	}
}

func CPUPercent(st container.StatsResponse) float64 {
	cpuDelta := float64(st.CPUStats.CPUUsage.TotalUsage) -
		float64(st.PreCPUStats.CPUUsage.TotalUsage)

	systemDelta := float64(st.CPUStats.SystemUsage) -
		float64(st.PreCPUStats.SystemUsage)

	if systemDelta == 0 || cpuDelta == 0 {
		return 0
	}

	numCPUs := float64(len(st.CPUStats.CPUUsage.PercpuUsage))
	if numCPUs == 0 {
		numCPUs = 1
	}

	return (cpuDelta / systemDelta) * numCPUs * 100.0
}

func MemoryPercent(st container.StatsResponse) float64 {
	used := float64(st.MemoryStats.Usage - st.MemoryStats.Stats["cache"])
	limit := float64(st.MemoryStats.Limit)

	if limit == 0 {
		return 0
	}

	return (used / limit) * 100.0
}
