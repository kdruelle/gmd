package containers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker"
)

type LoadedMsg struct {
	Containers []ContainerItem
	Err        error
}

type ContainerActionMsg struct {
	ContainerID string
	Action      string
	Err         error
}

func (m Model) FetchCmd() tea.Cmd {
	return func() tea.Msg {
		containers := m.client.Containers()
		containersItems := make([]ContainerItem, len(containers))
		for i, cont := range containers {
			containersItems[i] = NewContainerItem(cont)
		}
		return LoadedMsg{Containers: containersItems, Err: nil}
	}
}

func StartContainerCmd(clien *docker.Monitor, id string) tea.Cmd {
	return func() tea.Msg {
		msg := ContainerActionMsg{ContainerID: id, Action: "start"}
		msg.Err = clien.StartContainer(id)
		return msg
	}
}

func RestartContainerCmd(clien *docker.Monitor, id string) tea.Cmd {
	return func() tea.Msg {
		msg := ContainerActionMsg{ContainerID: id, Action: "restart"}
		msg.Err = clien.RestartContainer(id)
		return msg
	}
}
