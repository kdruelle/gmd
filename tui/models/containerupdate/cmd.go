package containerupdate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/container"
	"github.com/kdruelle/gmd/docker"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/tui/controllers/containerupdate"
)

type ContainerConfigMsg struct {
	Config *container.InspectResponse
}

type UpdateFinishedMsg struct {
}

func ContainerConfigCmd(client *docker.Monitor, id string) tea.Cmd {
	return func() tea.Msg {
		config, _ := client.GetContainerRawConfig(id)
		return ContainerConfigMsg{Config: config}
	}
}

func startUpdate(c *containerupdate.Controller, container client.Container) tea.Cmd {
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
