package images

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
	c, ok := item.(ImageItem)
	if !ok {
		return
	}

	title := style.Title().Render(c.Title())
	desc := style.Subtitle().Render(c.Description())

	content := lipgloss.JoinVertical(lipgloss.Left, title, desc)

	// On récupère quand même l'état sélectionné du delegate original
	if index == m.Index() {
		content = style.ListSelectedLine().Inherit(style.Bold()).Render(content)
		// title = d.Styles.SelectedTitle.Render(title)
		// desc = d.Styles.SelectedDesc.Render(desc)
	}
	// else {
	// 	title = d.Styles.NormalTitle.Render(title)
	// 	desc = d.Styles.NormalDesc.Render(desc)
	// }

	// Ta customisation ici : lipgloss partout, couleurs, emoji, flair…
	//_, _ = fmt.Fprintf(w, "%s\n%s", title, desc)
	fmt.Fprint(w, lipgloss.JoinHorizontal(lipgloss.Center, content, " " /*, c.statsContent*/))
}
