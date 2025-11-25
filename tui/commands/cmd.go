package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker"
)

func UpdateContainerCmd(c docker.Container) tea.Cmd {
	return func() tea.Msg {
		return UpdateContainerMsg{Container: c}
	}
}

func BackCmd() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

// func StartPullCmd(c chan PullProgressMsg) tea.Cmd {
// 	return func() tea.Msg {
// 		return PullStartedMsg{
// 			Channel: c,
// 		}
// 	}
// }

// func PullImageCmd(ch <-chan PullProgressMsg) tea.Cmd {
// 	return func() tea.Msg {
// 		msg, ok := <-ch
// 		if !ok {
// 			return PullCompleteMsg{}
// 		}
// 		return msg
// 	}
// }
