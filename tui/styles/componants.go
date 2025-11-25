package style

import "github.com/charmbracelet/lipgloss"

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorPrimary).
	Padding(0, 1)

var Subtitle = lipgloss.NewStyle().
	Foreground(ColorSecondary).
	Italic(true).
	PaddingLeft(2)

var Card = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorMuted).
	Background(ColorBg).
	Foreground(ColorFg).
	Padding(1, 2).
	Margin(1)

var Panel = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder()).
	BorderForeground(ColorSecondary).
	Padding(1)

var StatusBar = lipgloss.NewStyle().
	Background(ColorPrimary).
	Foreground(ColorBg).
	Bold(true).
	Padding(0, 1)

var SelectedItem = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
	Bold(true).
	Padding(0, 0, 0, 1)

var InactiveItem = lipgloss.NewStyle().
	Foreground(ColorMuted)

var ActiveItem = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSuccess)

var WarningItem = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorWarning)

var DangerItem = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorDanger)

var Separator = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorFg)

var Spinner = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorPrimary)
