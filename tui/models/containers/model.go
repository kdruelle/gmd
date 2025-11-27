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
	"github.com/kdruelle/gmd/docker/cache"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/controllers/containerstats"
	style "github.com/kdruelle/gmd/tui/styles"
)

type Model struct {
	cli                   *client.Client
	cache                 *cache.Cache
	list                  list.Model
	loaded                bool
	status                string
	all                   bool
	statsController       *containerstats.Controller
	checkUpdateInProgress map[string]struct{}
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

func New(cli *client.Client, cache *cache.Cache) Model {

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

	m := Model{
		cli:                   cli,
		cache:                 cache,
		list:                  l,
		all:                   false,
		checkUpdateInProgress: make(map[string]struct{}),
		//imgs:   images,
	}

	m.statsController = containerstats.New(cli)
	//m.statsController.Start()

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(WaitStatsEvent(m.statsController.Events()))
}

func (m Model) IsSearching() bool {
	return m.list.IsFiltered()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.toggleAll):
			m.ToggleAll()
			return m, nil

		case key.Matches(msg, keyMap.showLogs):
			cmd := exec.Command("docker", "logs", "-f", "--tail=200", m.list.SelectedItem().(ContainerItem).id)
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return nil
			})

		case key.Matches(msg, keyMap.restartContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && !slices.Contains([]string{container.StateRunning, container.StateRestarting}, c.state) {
				return m, RestartContainerCmd(m.cli, c.id)
			}
			return m, nil

		case key.Matches(msg, keyMap.startContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && !slices.Contains([]string{container.StateRunning, container.StateRestarting}, c.state) {
				return m, StartContainerCmd(m.cli, c.id)
			}
			return m, nil
		case key.Matches(msg, keyMap.updateContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok {
				if c.update != nil && *c.update {
					c, _ := m.cache.Container(c.id)
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
	case cache.Event:
		if msg.EventType == cache.ContainersLoadedEventType {
			if !m.loaded {
				m.loaded = true
			}
			log.Printf("received container event %+v", msg)
			cmds = append(cmds, m.initialLoad())
		}
		if msg.EventType == cache.ContainerEventType {
			if m.loaded {
				log.Printf("received container event %+v", msg)
				cmds = append(cmds, m.updateContainer(msg))
			}

		}
	case ContainerUpdateMsg:
		log.Printf("received container update event %+v", msg)
		if msg.Err == nil {
			for i, c := range m.list.Items() {
				if container, ok := c.(ContainerItem); ok && container.id == msg.ContainerID {
					b := msg.Update
					container.update = &b
					container.RenderContent()
					m.list.SetItem(i, container)
					break
				}
			}
		} else {
			log.Printf("error checking update for container %s: %s", msg.ContainerID, msg.Err)
		}
		delete(m.checkUpdateInProgress, msg.ContainerID)
	case ContainerActionMsg:
		if msg.Err != nil {
			m.status = style.DangerItem.Render(msg.Err.Error())
		}
	case containerstats.StatsMsg:
		for i, c := range m.list.Items() {
			if container, ok := c.(ContainerItem); ok && container.id == msg.ID {
				container.RenderStats(msg.Stats)
				m.list.SetItem(i, container)
				break
			}
		}
		return m, WaitStatsEvent(m.statsController.Events())
	}

	newList, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.list = newList
	return m, tea.Batch(cmds...)
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

func (m *Model) initialLoad() tea.Cmd {

	var containers = m.cache.Containers()

	var cmds = make([]tea.Cmd, 0, len(containers))

	slices.SortFunc(containers, func(a, b client.Container) int {
		return strings.Compare(a.Name, b.Name)
	})

	itemList := make([]list.Item, 0, len(containers))
	for _, item := range containers {
		container := NewContainerItem(item)
		if m.all {
			container.show = true
		} else {
			container.show = item.State.Running || item.State.Restarting
		}
		container.RenderContent()
		//m.statsController.AddContainer(container.id)
		itemList = append(itemList, container)
		m.checkUpdateInProgress[container.id] = struct{}{}
		cmds = append(cmds, CheckContainerUpdate(m.cli, container.id))
	}
	m.list.SetItems(itemList)
	return tea.Batch(cmds...)
}

func (m *Model) updateContainer(msg cache.Event) tea.Cmd {
	newContainer, err := m.cache.Container(msg.ActorID)

	for i, item := range m.list.Items() {
		if item.(ContainerItem).id == msg.ActorID {
			switch err {
			case cache.ErrContainerNotFound:
				//m.statsController.RemoveContainer(id)
				m.list.RemoveItem(i)
			case nil:
				var cmd tea.Cmd = nil
				container := NewContainerItem(newContainer)
				if item.(ContainerItem).update != nil {
					container.update = item.(ContainerItem).update
				} else {
					if _, ok := m.checkUpdateInProgress[container.id]; !ok {
						m.checkUpdateInProgress[container.id] = struct{}{}
						cmd = CheckContainerUpdate(m.cli, container.id)
					}
				}
				container.RenderContent()
				m.list.SetItem(i, container)
				return cmd
			}
			return nil
		}
	}
	container := NewContainerItem(newContainer)
	container.RenderContent()
	items := m.list.Items()
	items = append(items, container)
	slices.SortFunc(items, func(a, b list.Item) int {
		return strings.Compare(a.(ContainerItem).Name(), b.(ContainerItem).Name())
	})
	//m.statsController.AddContainer(container.id)
	m.list.SetItems(items)
	m.checkUpdateInProgress[container.id] = struct{}{}
	return CheckContainerUpdate(m.cli, container.id)
}

func (m *Model) ToggleAll() {
	m.all = !m.all

	for i, item := range m.list.Items() {
		if c, ok := item.(ContainerItem); ok {
			if m.all {
				c.show = true
				m.list.SetItem(i, c)
			} else {
				c.show = false
				if c.state == container.StateRunning || c.state == container.StateRestarting {
					c.show = true
				}
				m.list.SetItem(i, c)
			}
		}
	}
}
