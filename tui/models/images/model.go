package images

import (
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kdruelle/gmd/docker/cache"
	"github.com/kdruelle/gmd/docker/client"
	"github.com/kdruelle/gmd/docker/types"
	style "github.com/kdruelle/gmd/tui/styles"
)

type Model struct {
	cli    *client.Client
	cache  *cache.Cache
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

func New(cli *client.Client, cache *cache.Cache) Model {

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
		cli:   cli,
		cache: cache,
		list:  l,
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
			m.status = style.StatusBar().Render("Deleting image " + m.list.SelectedItem().(ImageItem).Title())
			return m, m.DeleteImagesCmd(m.list.SelectedItem().(ImageItem).ID)
		}

	case DeleteImageMsg:
		if msg.Err != nil {
			m.status = style.Danger().Render(msg.Err.Error())
		} else {
			m.status = style.Success().Render("Image supprimée")
		}
		m.applyFilter()
	case cache.Event:
		if msg.EventType == cache.ImagesLoadedEventType {
			if !m.loaded {
				m.loaded = true
			}
			log.Printf("received images loaded event: %+v", msg)
			m.applyFilter()
		}
		if msg.EventType == cache.ImageEventType {
			if m.loaded {
				log.Printf("received image event: %+v", msg)
				m.updateImage(msg.ActorID)
				// m.applyFilter()
			}

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

	var images []types.Image
	if m.unused {
		images = m.cache.ImagesUnused()
	} else {
		images = m.cache.Images()
	}

	slices.SortFunc(images, func(a, b types.Image) int {
		return strings.Compare(a.Tag(), b.Tag())
	})

	itemList := make([]list.Item, 0, len(images))
	for _, item := range images {

		itemList = append(itemList, ImageItem(item))

	}
	m.list.SetItems(itemList)
}

func (m *Model) updateImage(id string) {
	newImage, err := m.cache.Image(id)
	for i, item := range m.list.Items() {
		if item.(ImageItem).ID == id {
			switch err {
			case cache.ErrImageNotFound:
				m.list.RemoveItem(i)
			case nil:
				m.list.SetItem(i, ImageItem(newImage))
			}
			return
		}
	}
	items := m.list.Items()
	items = append(items, ImageItem(newImage))
	slices.SortFunc(items, func(a, b list.Item) int {
		return strings.Compare(a.(ImageItem).Title(), b.(ImageItem).Title())
	})
	m.list.SetItems(items)
}
