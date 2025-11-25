package images

import (
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kdruelle/gmd/docker"
	style "github.com/kdruelle/gmd/tui/styles"
)

type Model struct {
	client *docker.Monitor
	list   list.Model
	loaded bool
	unused bool
	status string
}

type listKeyMap struct {
	toggleUnused key.Binding
	delete       key.Binding
}

var keyMap = &listKeyMap{
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete selection"),
	),
	toggleUnused: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "toggle unused only"),
	),
}

func New(client *docker.Monitor) Model {

	items := []list.Item{}

	l := list.New(items, newItemDelegate(), 0, 0)
	l.Title = "Images"
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.delete,
			keyMap.toggleUnused,
		}
	}

	return Model{
		client: client,
		list:   l,
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

	case ImagesLoadedMsg:
		if msg.Err != nil {
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
		case key.Matches(msg, keyMap.toggleUnused):
			m.unused = !m.unused
			m.applyFilter()
			return m, nil
		case key.Matches(msg, keyMap.delete):
			m.status = style.StatusBar.Render("Deleting image " + m.list.SelectedItem().(ImageItem).Title())
			return m, m.DeleteImagesCmd(m.list.SelectedItem().(ImageItem).ID)
		}

	case DeleteImageMsg:
		if msg.Err != nil {
			m.status = style.DangerItem.Render(msg.Err.Error())
		} else {
			m.status = style.ActiveItem.Render("Image supprimée")
		}
		m.applyFilter()
	case docker.Event:
		if msg.EventType == docker.ImageEventType {
			if !m.loaded {
				m.loaded = true
			}
			log.Printf("received image event: %v", msg)
			m.applyFilter()
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m Model) View() string {
	if !m.loaded {
		return "Chargement des images Docker..."
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
		m.status,
	)
}

func (m *Model) applyFilter() {

	log.Printf("image applying filter")

	var images []docker.Image
	if m.unused {
		images = m.client.ImagesUnused()
	} else {
		images = m.client.Images()
	}

	slices.SortFunc(images, func(a, b docker.Image) int {
		return strings.Compare(a.Tag(), b.Tag())
	})

	itemList := make([]list.Item, 0, len(images))
	for _, item := range images {

		itemList = append(itemList, ImageItem(item))

	}
	m.list.SetItems(itemList)
}
