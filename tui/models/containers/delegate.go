package containers

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	style "github.com/kdruelle/gmd/tui/styles"
)

var (
	col1Style = lipgloss.NewStyle().Width(30)
	col2Style = lipgloss.NewStyle().Width(20)
	col3Style = lipgloss.NewStyle().Width(90)
)

type ItemDelegate struct {
	list.DefaultDelegate
}

func newItemDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	return ItemDelegate{d}
}

func (d ItemDelegate) Height() int  { return 2 }
func (d ItemDelegate) Spacing() int { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	c, ok := item.(ContainerItem)
	if !ok {
		return
	}

	content := c.content
	if index == m.Index() {
		content = style.SelectedItem.Render(content)
	}

	fmt.Fprint(w, lipgloss.JoinHorizontal(lipgloss.Center, content, " ", c.statsContent))
}
