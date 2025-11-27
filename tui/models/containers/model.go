package containers

import (
	"errors"
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
	"github.com/kdruelle/gmd/docker/types"
	"github.com/kdruelle/gmd/tui/commands"
	"github.com/kdruelle/gmd/tui/controllers/containerstats"
	"github.com/kdruelle/gmd/tui/models/containerupdate"
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
	stopContainer    key.Binding
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
	stopContainer: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "stop container"),
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
			keyMap.stopContainer,
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
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && c.state == container.StateRunning {
				m.updateContainerActionState(c.id, container.StateRestarting)
				m.status = style.StatusBar.Render("Restarting container " + m.list.SelectedItem().(ContainerItem).name)
				return m, commands.ContainerCmd(m.cli, commands.RestartContainerAction, c.id)
			}
			return m, nil

		case key.Matches(msg, keyMap.startContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok && !slices.Contains([]string{container.StateRunning, container.StateRestarting}, c.state) {
				return m, commands.ContainerCmd(m.cli, commands.StartContainerAction, c.id)
			}
			return m, nil

		case key.Matches(msg, keyMap.stopContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok {
				return m, commands.ContainerCmd(m.cli, commands.StopContainerAction, c.id)
			}
			return m, nil

		case key.Matches(msg, keyMap.updateContainer):
			if c, ok := m.list.SelectedItem().(ContainerItem); ok {
				if c.update != nil && *c.update {
					c, _ := m.cache.Container(c.id)
					return m, commands.SwitchPageCmd(func() tea.Model {
						u := containerupdate.New(c, m.cli)
						return u
					})
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
				cmds = append(cmds, m.handleContainerEvent(msg))
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
	case commands.ContainerActionMsg:
		if msg.Err != nil {
			m.status = style.DangerItem.Render(msg.Err.Error())
		} else {
			m.status = ""
			m.updateContainerActionState(msg.ContainerID, "")
		}
		// case containerstats.StatsMsg:
		// 	for i, c := range m.list.Items() {
		// 		if container, ok := c.(ContainerItem); ok && container.id == msg.ID {
		// 			container.RenderStats(msg.Stats)
		// 			m.list.SetItem(i, container)
		// 			break
		// 		}
		// 	}
		// 	return m, WaitStatsEvent(m.statsController.Events())
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

// initialLoad loads all containers from the cache and sets the list with the
// retrieved containers. It also checks for updates for all containers
func (m *Model) initialLoad() tea.Cmd {

	var containers = m.cache.Containers()

	var cmds = make([]tea.Cmd, 0, len(containers))

	slices.SortFunc(containers, func(a, b types.Container) int {
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

// handleContainerEvent handles a container event from the cache.
//
// The function first retrieves the ContainerItem from the cache with the given id.
// If the container is not found, it removes the container from the list.
// If the container is found, it gets the old ContainerItem from the list with the given id.
// If the old container is not found, it adds the new container to the list.
// If the old container is found, it updates the old container with the new container.
//
// The function returns a tea.Cmd that executes the update if needed.
func (m *Model) handleContainerEvent(msg cache.Event) tea.Cmd {
	newContainer, err := m.cache.Container(msg.ActorID)

	if err != nil {
		m.removeContainer(msg.ActorID)
		return nil
	}

	oldContainer, oldContainerIdx, err := m.getContainerWithIndex(msg.ActorID)

	if err != nil {
		return m.addNewContainer(newContainer)
	}

	return m.updateContainer(newContainer, oldContainer, oldContainerIdx)
}

// getContainerWithIndex returns the ContainerItem with the given id and its index in the model's list.
// If the container is not found, it returns an empty ContainerItem, -1 as the index, and an error.
func (m *Model) getContainerWithIndex(id string) (ContainerItem, int, error) {
	for i, item := range m.list.Items() {
		if item.(ContainerItem).id == id {
			return item.(ContainerItem), i, nil
		}
	}
	return ContainerItem{}, -1, errors.New("container not found")
}

// removeContainer removes a container from the list.
//
// The function iterates over the list of containers and removes the first container that matches the given id.
// If no container is found, the function does nothing.
// The function does not return anything.
func (m *Model) removeContainer(id string) {
	for i, item := range m.list.Items() {
		if item.(ContainerItem).id == id {
			m.list.RemoveItem(i)
			return
		}
	}
}

// addNewContainer adds a new container to the list and checks for update.
//
// The function sorts the list by container name and adds the new container at the end.
// It also adds the container to the list of containers to check for update.
// The function returns a command to check for container update.
//
// The function also renders the content of the new container.
func (m *Model) addNewContainer(container types.Container) tea.Cmd {
	newContainer := NewContainerItem(container)
	newContainer.RenderContent()
	items := m.list.Items()
	items = append(items, newContainer)
	slices.SortFunc(items, func(a, b list.Item) int {
		return strings.Compare(a.(ContainerItem).Name(), b.(ContainerItem).Name())
	})
	m.list.SetItems(items)
	m.checkUpdateInProgress[newContainer.id] = struct{}{}
	return CheckContainerUpdate(m.cli, newContainer.id)
}

func (m *Model) updateContainer(newContainer types.Container, oldContainer ContainerItem, index int) tea.Cmd {
	var cmd tea.Cmd = nil
	c := NewContainerItem(newContainer)

	// update flag
	if oldContainer.update != nil {
		c.update = oldContainer.update
	} else {
		if _, ok := m.checkUpdateInProgress[c.id]; !ok {
			m.checkUpdateInProgress[c.id] = struct{}{}
			cmd = CheckContainerUpdate(m.cli, c.id)
		}
	}

	c.actionState = oldContainer.actionState

	c.RenderContent()
	m.list.SetItem(index, c)
	return cmd

}

func (m *Model) updateContainerActionState(id string, state string) {
	for i, item := range m.list.Items() {
		if c := item.(ContainerItem); c.id == id {
			c.actionState = state
			c.RenderContent()
			m.list.SetItem(i, c)
		}
	}
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
