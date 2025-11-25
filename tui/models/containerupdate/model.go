package containerupdate

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kdruelle/gmd/docker"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/controllers/containerupdate"
)

// type layer struct {
// 	Status      string
// 	ProgressMsg commands.PullProgressMsg
// 	Progress    progress.Model
// }

type Model struct {
	container docker.Container
	//layers        map[string]layer
	// progressOrder []string
	// done          bool
	client     *docker.Monitor
	updatech   chan containerupdate.ControllerUpdateMsg
	controller *containerupdate.Controller
	//progressChan  chan commands.PullProgressMsg

	// stopMessage  string
	// errorMessage string
}

type listKeyMap struct {
	returnKey key.Binding
}

var keyMap = &listKeyMap{
	returnKey: key.NewBinding(
		key.WithKeys("esc", "enter"),
		key.WithHelp("enter", "get back to main menu"),
	),
}

func New(c docker.Container, client *docker.Monitor) Model {
	uch := make(chan containerupdate.ControllerUpdateMsg)
	controller := containerupdate.New(client, uch)
	return Model{
		container:  c,
		client:     client,
		updatech:   uch,
		controller: controller,

		//layers:    make(map[string]layer),
	}
}

func (m Model) Init() tea.Cmd {
	log.Printf("init update for container %s", m.container.Name)
	// ch := startPull(m.client, m.container.Config.Image)
	return startUpdate(m.controller, m.container)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case containerupdate.ControllerUpdateMsg:
		log.Printf("received update msg: %+v", msg)
		_ = msg
		return m, waitUpdateEvent(m.updatech)
	case UpdateFinishedMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.returnKey):
			return m, commands.BackCmd()
		}
	// case commands.PullStartedMsg:
	// 	m.progressChan = msg.Channel
	// 	return m, commands.PullImageCmd(m.progressChan)
	// case commands.PullProgressMsg:
	// 	if msg.Err != nil {
	// 		m.errorMessage = msg.Err.Error()
	// 		return m, nil
	// 	}
	// 	var p progress.Model
	// 	if _, exists := m.layers[msg.LayerID]; !exists {
	// 		m.progressOrder = append(m.progressOrder, msg.LayerID) // première fois qu’on voit ce layer
	// 		p = progress.New(progress.WithDefaultGradient(), progress.WithWidth(40))
	// 		m.layers[msg.LayerID] = layer{Status: msg.Status, Progress: p, ProgressMsg: msg}
	// 	}

	// 	layer := m.layers[msg.LayerID]
	// 	layer.Status = msg.Status
	// 	layer.ProgressMsg = msg

	// 	if msg.Status == "Pull complete" || msg.Status == "Already exists" || msg.Status == "Verifying" || msg.Status == "Download complete" {
	// 		msg.ProgressPct = 1
	// 	}
	// 	log.Printf("set pct to %f", msg.ProgressPct)
	// 	cmd := layer.Progress.SetPercent(msg.ProgressPct)

	// 	m.layers[msg.LayerID] = layer
	// 	return m, tea.Batch(commands.PullImageCmd(m.progressChan), cmd)
	// case commands.PullCompleteMsg:
	// 	m.stopMessage = fmt.Sprintf("⠋ Stopping container %s", m.container.Name)
	// 	return m, ContainerConfigCmd(m.client, m.container.ID)
	case ContainerConfigMsg:
		// stop
		// case progress.FrameMsg:
		// 	for i, layer := range m.layers {
		// 		np, cmd := layer.Progress.Update(msg)
		// 		layer.Progress = np.(progress.Model)
		// 		m.layers[i] = layer
		// 		return m, cmd
		// 	}
	}
	return m, nil
}

func (m Model) View() string {
	var rows []string

	header := lipgloss.NewStyle().Bold(true).Render(
		fmt.Sprintf("Updating container %s ...", m.container.Name),
	)

	rows = append(rows, header)
	rows = append(rows, m.controller.GetLines()...)

	// Format 1 ligne par layer avec un vrai layout propre
	// for i, id := range m.progressOrder {
	// 	layer := m.layers[id]

	// 	if i == 0 {
	// 		rows = append(rows, layer.Status)
	// 		continue
	// 	}

	// 	status := lipgloss.NewStyle().
	// 		Width(20).
	// 		Foreground(lipgloss.Color("#88C0D0")).
	// 		Render(layer.Status)

	// 	bar := layer.ProgressMsg.Progress
	// 	// bar := layer.Progress.View()

	// 	row := lipgloss.JoinHorizontal(lipgloss.Top, status, bar)
	// 	rows = append(rows, row)
	// }

	// erreurs
	// if m.errorMessage != "" {
	// 	errBox := lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#BF616A")).
	// 		Render(m.errorMessage)
	// 	rows = append(rows, "", errBox)
	// }

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
