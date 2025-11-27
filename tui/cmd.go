package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker/cache"
)

type CacheStartMsg struct {
	Err error
}

func StartMonitorCache(m *cache.Cache) tea.Cmd {
	return func() tea.Msg {
		err := m.LoadAndStart()
		return CacheStartMsg{Err: err}
	}
}

func WaitDockerEvent(ch <-chan cache.Event) tea.Cmd {
	return func() tea.Msg {
		//var now = time.Now()
		var e cache.Event

		for {
			e = <-ch
			//if e.EventType != cache.ContainerStatsEventType || time.Since(now) > 1*time.Second {
			return e
			//}
		}
	}
}

func SendResize(width, height int) tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		}
	}
}
