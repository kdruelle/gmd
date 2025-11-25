package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker"
)

type MonitorStartMsg struct {
}

func StartMonitor(m *docker.Monitor) tea.Cmd {
	return func() tea.Msg {
		m.Start()
		return MonitorStartMsg{}
	}
}

func WaitDockerEvent(ch <-chan docker.Event) tea.Cmd {
	return func() tea.Msg {
		var now = time.Now()
		var e docker.Event

		for {
			e = <-ch
			if e.EventType != docker.ContainerStatsEventType || time.Since(now) > 1*time.Second {
				return e
			}
		}
	}
}
