package containerupdate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker/types"
	"github.com/kdruelle/gmd/tui/controllers/containerupdate"
)

type UpdateFinishedMsg struct {
}

func startUpdate(c *containerupdate.Controller, container types.Container) tea.Cmd {
	return func() tea.Msg {
		c.StartUpdate(container)
		return containerupdate.ControllerUpdateMsg{}
	}
}

func waitUpdateEvent(updatech <-chan containerupdate.ControllerUpdateMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-updatech
		if !ok {
			return UpdateFinishedMsg{}
		}
		return msg
	}
}
