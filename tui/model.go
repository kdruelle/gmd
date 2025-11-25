package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/componants"
	"github.com/kdruelle/gmd/tui/models/containerupdate"
	"github.com/kdruelle/gmd/tui/models/maintab"
)

// ---------------------------------------------------
// Model Root
// ---------------------------------------------------

type Model struct {
	client *docker.Monitor
	stack  []tea.Model
}

func NewModel() Model {
	client, _ := docker.NewMonitor()

	mainModel := maintab.New(client)

	m := Model{
		client: client,
	}

	m.stack = []tea.Model{
		mainModel,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	top := m.stack[len(m.stack)-1]
	return tea.Batch(
		StartMonitor(m.client),
		WaitDockerEvent(m.client.Events()),
		top.Init(),
	)
}

// ---------------------------------------------------
// Update
// ---------------------------------------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		if searchable, ok := m.stack[len(m.stack)-1].(componants.Searchable); ok && searchable.IsSearching() {
			var cmd tea.Cmd
			m.stack[len(m.stack)-1], cmd = m.stack[len(m.stack)-1].Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case docker.Event:
		var cmd tea.Cmd
		m.stack[0], cmd = m.stack[0].Update(msg)

		return m, tea.Batch(WaitDockerEvent(m.client.Events()), cmd)
	case commands.UpdateContainerMsg:
		u := containerupdate.New(msg.Container, m.client)
		cmd := u.Init()
		m.stack = append(m.stack, u)
		return m, tea.Batch(tea.ExitAltScreen, cmd)
	case commands.BackMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1] // pop
		}
		return m, tea.EnterAltScreen
	}

	top := m.stack[len(m.stack)-1]
	newTop, cmd := top.Update(msg)
	m.stack[len(m.stack)-1] = newTop

	return m, cmd
}

// ---------------------------------------------------
// View
// ---------------------------------------------------

func (m Model) View() string {
	top := m.stack[len(m.stack)-1]
	return top.View()
}
