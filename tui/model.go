package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kdruelle/gmd/docker/cache"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/componants"
	"github.com/kdruelle/gmd/tui/models/containers"
	"github.com/kdruelle/gmd/tui/models/maintab"
)

// ---------------------------------------------------
// Model Root
// ---------------------------------------------------

type Model struct {
	cli         *client.Client
	dockerCache *cache.Cache
	stack       []tea.Model

	screeWidth   int
	screenHeight int
}

func NewModel() Model {
	cli, _ := client.NewClient()
	cache := cache.NewCache(cli)

	mainModel := maintab.New(cli, cache)

	m := Model{
		cli:         cli,
		dockerCache: cache,
	}

	m.stack = []tea.Model{
		mainModel,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	top := m.stack[len(m.stack)-1]
	return tea.Batch(
		StartMonitorCache(m.dockerCache),
		WaitDockerEvent(m.dockerCache.Events()),
		top.Init(),
	)
}

// ---------------------------------------------------
// Update
// ---------------------------------------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.screeWidth = msg.Width
		m.screenHeight = msg.Height

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

	case cache.Event:
		var cmd tea.Cmd
		m.stack[0], cmd = m.stack[0].Update(msg)
		return m, tea.Batch(WaitDockerEvent(m.dockerCache.Events()), cmd)

	case containers.ContainerUpdateMsg:
		var cmd tea.Cmd
		m.stack[0], cmd = m.stack[0].Update(msg)
		return m, cmd
		
	case commands.SwitchPageMsg:
		model := msg.Model
		if model == nil {
			m.stack = m.stack[:len(m.stack)-1] // pop
			return m, nil
		}
		cmd := model.Init()
		m.stack = append(m.stack, model)
		return m, tea.Batch( /*tea.ExitAltScreen,*/ cmd, SendResize(m.screeWidth, m.screenHeight))
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
