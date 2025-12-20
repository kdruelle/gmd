package containers

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	style "github.com/kdruelle/gmd/tui/styles"
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
	//prefix := " "
	if index == m.Index() {
		content = c.Render(true)
		content = style.ListSelectedLine().Inherit(style.Bold()).Render(content)
		// content = lipgloss.JoinHorizontal(
		// 	0,
		// 	style.SelectedBar.Render("▍ "),
		// 	style.SelectedLine.Render(content),
		// )
		//prefix = style.SelectedBar.Render(lipgloss.JoinVertical(lipgloss.Center, "▍ ", "▍ "))
	}

	fmt.Fprint(w, lipgloss.JoinHorizontal(lipgloss.Center, content, " " /*, c.statsContent*/))
}
