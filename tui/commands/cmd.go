package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker/client"
)

func SwitchPageCmd(modelCreate func() tea.Model) tea.Cmd {
	return func() tea.Msg {
		var model tea.Model = nil
		if modelCreate != nil {
			model = modelCreate()
		}
		return SwitchPageMsg{Model: model}
	}
}

func ContainerCmd(cli *client.Client, action Action, id string) tea.Cmd {
	return func() tea.Msg {
		msg := ContainerActionMsg{ContainerID: id, Action: action}

		switch action {
		case StartContainerAction:
			msg.Err = cli.StartContainer(id)
		case StopContainerAction:
			msg.Err = cli.StopContainer(id)
		case RestartContainerAction:
			msg.Err = cli.RestartContainer(id)
		}
		return msg
	}
}
