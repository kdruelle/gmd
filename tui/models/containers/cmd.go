package containers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/tui/controllers/containerstats"
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

type ContainerUpdateMsg struct {
	ContainerID string
	Update      bool
	Err         error
}

func (m Model) FetchCmd() tea.Cmd {
	return func() tea.Msg {
		containers := m.cache.Containers()
		containersItems := make([]ContainerItem, len(containers))
		for i, cont := range containers {
			containersItems[i] = NewContainerItem(cont)
		}
		return LoadedMsg{Containers: containersItems, Err: nil}
	}
}

func StartContainerCmd(cli *client.Client, id string) tea.Cmd {
	return func() tea.Msg {
		msg := ContainerActionMsg{ContainerID: id, Action: "start"}
		msg.Err = cli.StartContainer(id)
		return msg
	}
}

func RestartContainerCmd(cli *client.Client, id string) tea.Cmd {
	return func() tea.Msg {
		msg := ContainerActionMsg{ContainerID: id, Action: "restart"}
		msg.Err = cli.RestartContainer(id)
		return msg
	}
}

func CheckContainerUpdate(cli *client.Client, id string) tea.Cmd {
	return func() tea.Msg {
		update, err := cli.CheckUpdate(id)
		return ContainerUpdateMsg{ContainerID: id, Update: update, Err: err}
	}
}

func WaitStatsEvent(ch <-chan containerstats.StatsMsg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
