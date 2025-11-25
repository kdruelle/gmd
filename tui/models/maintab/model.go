package maintab

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kdruelle/gmd/docker"
	"github.com/kdruelle/gmd/tui/componants"
	"github.com/kdruelle/gmd/tui/models/containers"
	"github.com/kdruelle/gmd/tui/models/images"
	style "github.com/kdruelle/gmd/tui/styles"
)

// ---------------------------------------------------
// Messages
// ---------------------------------------------------

type SwitchTabMsg int

// ---------------------------------------------------
// Model Root
// ---------------------------------------------------

const (
	imagesTabIndex     = 0
	containersTabIndex = 1
)

type Model struct {
	client    *docker.Monitor
	lists     []componants.ListModel
	activeTab int
}

func New(client *docker.Monitor) Model {

	m := Model{
		client: client,
		lists:  make([]componants.ListModel, 2),
	}

	m.lists[imagesTabIndex] = images.New(client)
	m.lists[containersTabIndex] = containers.New(client)
	return m
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)
	for i := range m.lists {
		cmds = append(cmds, m.lists[i].Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) IsSearching() bool {
	return m.lists[m.activeTab].(componants.Searchable).IsSearching()
}

// ---------------------------------------------------
// Update
// ---------------------------------------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "tab", "ctrl+tab":
			// On avance d’un onglet, circulation circulaire
			m.activeTab = (m.activeTab + 1) % len(m.lists)
			return m, nil

		case "shift+tab":
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = len(m.lists) - 1
			}
			return m, nil

		}

		// Pass key stroke to active tab
		l, cmd := m.lists[m.activeTab].Update(msg)
		m.lists[m.activeTab] = l
		return m, cmd

	case docker.Event:
		switch msg.EventType {
		case docker.ImageEventType:
			l, cmd := m.lists[imagesTabIndex].Update(msg)
			m.lists[imagesTabIndex] = l
			return m, cmd
		case docker.ContainerEventType, docker.ContainerStatsEventType:
			l, cmd := m.lists[containersTabIndex].Update(msg)
			m.lists[containersTabIndex] = l
			return m, cmd
		}
	}

	// pass all events to all lists
	var cmds []tea.Cmd
	for i := range m.lists {
		l, cmd := m.lists[i].Update(msg)
		cmds = append(cmds, cmd)
		m.lists[i] = l
	}

	return m, tea.Batch(cmds...)
}

// ---------------------------------------------------
// View
// ---------------------------------------------------

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewTabs(),
		style.Separator.Render("────────────────────────────────────────────────────────────"),
		m.viewContent(),
	)
}

func (m Model) viewTabs() string {

	var (
		tabImages     = style.InactiveItem.Render(" Images ")
		tabContainers = style.InactiveItem.Render(" Containers ")
	)

	switch m.activeTab {
	case imagesTabIndex:
		tabImages = style.ActiveItem.Render(" Images ")
	case containersTabIndex:
		tabContainers = style.ActiveItem.Render(" Containers ")
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, tabImages, tabContainers)
}

func (m Model) viewContent() string {
	return m.lists[m.activeTab].View()
}
