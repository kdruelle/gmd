package style

import "github.com/charmbracelet/lipgloss"

func Title() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(current.Primary).PaddingLeft(2)
}

func Normal() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(current.Fg)
}
func Subtitle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(current.Secondary).Italic(true).PaddingLeft(2)
}

func StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(current.Primary).
		Foreground(current.Bg).
		Bold(true).
		Padding(0, 1)
}

func ListSelectedLine() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Padding(0, 0, 0, 0)
}

func Spinner() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(current.Primary).Bold(true)
}
// var Details = lipgloss.NewStyle().
// 	Foreground(ColorSecondary).
// 	PaddingLeft(2)

// var Bold = lipgloss.NewStyle().Bold(true)

// var Card = lipgloss.NewStyle().
// 	Border(lipgloss.RoundedBorder()).
// 	BorderForeground(ColorMuted).
// 	Background(ColorBg).
// 	Foreground(ColorFg).
// 	Padding(1, 2).
// 	Margin(1)

// var Panel = lipgloss.NewStyle().
// 	Border(lipgloss.ThickBorder()).
// 	BorderForeground(ColorSecondary).
// 	Padding(1)

// var SelectedItem = lipgloss.NewStyle().
// 	Border(lipgloss.NormalBorder(), false, false, true, true).
// 	BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
// 	Bold(true).
// 	Padding(0, 0, 0, 0)

// var SelectedBar = lipgloss.NewStyle().
// 	Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

// var ListSelectedLine = lipgloss.NewStyle().
// 	Border(lipgloss.ThickBorder(), false, false, false, true).
// 	BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
// 	Padding(0, 0, 0, 0)

