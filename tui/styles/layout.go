package style

import "github.com/charmbracelet/lipgloss"

var (
    ColumnLeft = lipgloss.NewStyle().
        Width(30).
        PaddingRight(1)

    ColumnRight = lipgloss.NewStyle().
        Width(80).
        PaddingLeft(1)

    FullWidth = lipgloss.NewStyle().
        Width(110)

    Centered = lipgloss.NewStyle().
        Align(lipgloss.Center)
)
