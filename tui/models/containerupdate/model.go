package containerupdate

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/controllers/containerupdate"
)

type Model struct {
	container  client.Container
	cli        *client.Client
	controller *containerupdate.Controller
	screenW    int
	screenH    int

	titleBlock string
	completed  bool
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

func New(c client.Container, client *client.Client) Model {
	controller := containerupdate.New(client)
	m := Model{
		container:  c,
		cli:        client,
		controller: controller,
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Width(90).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Updating container %s ...", m.container.Name))

	m.titleBlock = lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
	)

	return m
}

func (m Model) Init() tea.Cmd {
	log.Printf("init update for container %s", m.container.Name)
	return startUpdate(m.controller, m.container)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenW = msg.Width
		m.screenH = msg.Height
	case containerupdate.ControllerUpdateMsg:
		_ = msg
		return m, waitUpdateEvent(m.controller.Events())
	case UpdateFinishedMsg:
		m.completed = true
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.returnKey):
			if m.completed {
				return m, commands.BackCmd()
			}
		}

	}
	return m, nil
}

func (m Model) View() string {

	contentLines := lipgloss.JoinVertical(
		lipgloss.Left,
		m.controller.GetLines()...,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.titleBlock,
		contentLines,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#81A1C1")).
		Padding(1, 2).
		Width(90).
		Height(25).
		Align(lipgloss.Left)

	return lipgloss.Place(
		m.screenW, m.screenH,
		lipgloss.Center, lipgloss.Center,
		box.Render(content),
	)
}
