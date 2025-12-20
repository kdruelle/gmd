package style

import "github.com/charmbracelet/lipgloss"

func Bold() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

func Italic() lipgloss.Style {
	return lipgloss.NewStyle().Italic(true)
}

func Danger() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorDanger())
}

func Warning() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorWarning())
}

func Success() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorSuccess())
}

func Inactive() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorInactive())
}
