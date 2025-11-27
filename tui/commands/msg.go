package commands

import tea "github.com/charmbracelet/bubbletea"

type SwitchPageMsg struct {
	Model tea.Model
}

type Action string

const (
	StartContainerAction   Action = "start"
	StopContainerAction    Action = "stop"
	RestartContainerAction Action = "restart"
)

type ContainerActionMsg struct {
	ContainerID string
	Action      Action
	Update      bool
	Err         error
}

// type PullStartedMsg struct {
// 	Channel chan PullProgressMsg
// }

// type PullProgressMsg struct {
// 	LayerID         string
// 	Status          string
// 	Progress        string
// 	ProgressCurrent float64
// 	ProgressTotal   float64
// 	ProgressPct     float64
// 	Err             error
// }

// type PullCompleteMsg struct {
// }

// type StoppedContainerMsg struct {
// 	Err error
// }
