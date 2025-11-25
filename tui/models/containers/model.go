package containers

import (
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/kdruelle/gmd/docker"
	"github.com/kdruelle/gmd/tui/commands"
	style "github.com/kdruelle/gmd/tui/styles"
)

type Model struct {
	client *docker.Monitor
	list   list.Model
	loaded bool
	status string
	all    bool
}

type listKeyMap struct {
	toggleAll        key.Binding
	showLogs         key.Binding
	restartContainer key.Binding
	startContainer   key.Binding
	updateContainer  key.Binding
	execTerminal     key.Binding
}

var keyMap = &listKeyMap{
	toggleAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "toggle all containers"),
	),
	showLogs: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "show logs"),
	),
	restartContainer: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restart container"),
	),
	startContainer: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start container"),
	),
	updateContainer: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "update container"),
	),
	execTerminal: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "open terminal"),
	),
}

func New(client *docker.Monitor) Model {

	items := []list.Item{}

	l := list.New(items, newItemDelegate(), 0, 0)
	l.Title = "Containers"
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.toggleAll,
			keyMap.showLogs,
			keyMap.updateContainer,
			keyMap.restartContainer,
			keyMap.startContainer,
			keyMap.execTerminal,
		}
	}

	return Model{
		client: client,
		list:   l,
		all:    false,
		//imgs:   images,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) IsSearching() bool {
	return m.list.IsFiltered()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case LoadedMsg:
		if msg.Err != nil {
			m.status = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(msg.Err.Error())
			return m, nil
		}
		m.loaded = true
		m.applyFilter()
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.toggleAll):
			m.all = !m.all
			m.applyFilter()
			return m, nil

		case key.Matches(msg, keyMap.showLogs):
			cmd := exec.Command("docker", "logs", "-f", "--tail=200", m.list.SelectedItem().(ContainerItem).id)
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return nil
			})

		case key.Matches(msg, keyMap.restartContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && !slices.Contains([]string{container.StateRunning, container.StateRestarting}, c.state) {
				return m, RestartContainerCmd(m.client, c.id)
			}
			return m, nil

		case key.Matches(msg, keyMap.startContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && !slices.Contains([]string{container.StateRunning, container.StateRestarting}, c.state) {
				return m, StartContainerCmd(m.client, c.id)
			}
			return m, nil
		case key.Matches(msg, keyMap.updateContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok {
				if c.update {
					c, _ := m.client.Container(c.id)
					return m, commands.UpdateContainerCmd(c)
				}
			}
			return m, nil
		case key.Matches(msg, keyMap.execTerminal):
			cmd := exec.Command("docker", "exec", "-it", m.list.SelectedItem().(ContainerItem).id, "/bin/sh")
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return nil
			})
		}
	case docker.Event:
		if msg.EventType == docker.ContainerEventType {
			if !m.loaded {
				m.loaded = true
			}
			log.Printf("received container event %+v", msg)
			m.applyFilter()
		} else if msg.EventType == docker.ContainerStatsEventType {
			if !m.loaded {
				return m, nil
			}
			log.Printf("received container stats event %+v", msg)
			m.applyFilter()
		}
	case ContainerActionMsg:
		if msg.Err != nil {
			m.status = style.DangerItem.Render(msg.Err.Error())
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m Model) View() string {
	if !m.loaded {
		return "Chargement des containers Docker..."
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
		m.status,
	)
}

func (m *Model) applyFilter() {

	var containers = m.client.Containers()

	if !m.all {
		var items = make([]docker.Container, 0, len(containers))
		for _, item := range containers {
			if item.State.Running || item.State.Restarting {
				items = append(items, item)
			}
		}
		containers = items
	}

	slices.SortFunc(containers, func(a, b docker.Container) int {
		return strings.Compare(a.Name, b.Name)
	})
	// if m.unused {
	// 	images = m.client.ImagesUnused()
	// } else {
	// 	images = m.client.Images()
	// }

	itemList := make([]list.Item, 0, len(containers))
	for _, item := range containers {

		itemList = append(itemList, NewContainerItem(item))

	}
	m.list.SetItems(itemList)
}

func (m *Model) updateContainer(id string) {
	newContainer, err := m.client.Container(id)
	for i, item := range m.list.Items() {
		if item.(ContainerItem).id == id {
			switch err {
			case docker.ErrContainerNotFound:
				m.list.RemoveItem(i)
			case nil:
				m.list.SetItem(i, NewContainerItem(newContainer))
			}
			return
		}
	}
	items := m.list.Items()
	items = append(items, NewContainerItem(newContainer))
	slices.SortFunc(items, func(a, b list.Item) int {
		return strings.Compare(a.(ContainerItem).Name(), b.(ContainerItem).Name())
	})
	m.list.SetItems(items)
}
