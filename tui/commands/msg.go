package commands

import "github.com/kdruelle/gmd/docker"

type BackMsg struct{}

type UpdateContainerMsg struct {
	Container docker.Container
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
