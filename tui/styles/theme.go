package style

import "github.com/charmbracelet/lipgloss"

var current Theme = Default

type Theme struct {
	Bg        lipgloss.AdaptiveColor
	Fg        lipgloss.AdaptiveColor
	Muted     lipgloss.AdaptiveColor
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Success   lipgloss.AdaptiveColor
	Warning   lipgloss.AdaptiveColor
	Danger    lipgloss.AdaptiveColor
}

func ColorSuccess() lipgloss.AdaptiveColor  { return current.Success }
func ColorWarning() lipgloss.AdaptiveColor  { return current.Warning }
func ColorDanger() lipgloss.AdaptiveColor   { return current.Danger }
func ColorInactive() lipgloss.AdaptiveColor { return current.Muted }

var Default = Theme{
	Bg:        lipgloss.AdaptiveColor{Light: "#F6F8FA", Dark: "#2E3440"},
	Fg:        lipgloss.AdaptiveColor{Light: "#24292E", Dark: "#ECEFF4"},
	Muted:     lipgloss.AdaptiveColor{Light: "#6A737D", Dark: "#4C566A"},
	Primary:   lipgloss.AdaptiveColor{Light: "#0366D6", Dark: "#88C0D0"},
	Secondary: lipgloss.AdaptiveColor{Light: "#6A737D", Dark: "#81A1C1"},
	Success:   lipgloss.AdaptiveColor{Light: "#22863A", Dark: "#A3BE8C"},
	Warning:   lipgloss.AdaptiveColor{Light: "#D29922", Dark: "#EBCB8B"},
	Danger:    lipgloss.AdaptiveColor{Light: "#B31D28", Dark: "#BF616A"},
}

var NordHub = Theme{
	Bg:        lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#2E3440"},
	Fg:        lipgloss.AdaptiveColor{Light: "#1C2128", Dark: "#ECEFF4"},
	Muted:     lipgloss.AdaptiveColor{Light: "#6E7781", Dark: "#4C566A"},
	Primary:   lipgloss.AdaptiveColor{Light: "#0366D6", Dark: "#88C0D0"},
	Secondary: lipgloss.AdaptiveColor{Light: "#586069", Dark: "#81A1C1"},
	Success:   lipgloss.AdaptiveColor{Light: "#22863A", Dark: "#A3BE8C"},
	Warning:   lipgloss.AdaptiveColor{Light: "#C18310", Dark: "#EBCB8B"},
	Danger:    lipgloss.AdaptiveColor{Light: "#B31D28", Dark: "#BF616A"},
}
