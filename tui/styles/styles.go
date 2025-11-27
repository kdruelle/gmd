package style

import "github.com/charmbracelet/lipgloss"

var (
	// Palette principale
	ColorPrimary   lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#0366D6", Dark: "#88C0D0"}
	ColorSecondary lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#6A737D", Dark: "#81A1C1"}
	ColorSuccess   lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#22863A", Dark: "#A3BE8C"}
	ColorDanger    lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#B31D28", Dark: "#BF616A"}
	ColorWarning   lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#D29922", Dark: "#EBCB8B"}
	ColorBg        lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#F6F8FA", Dark: "#2E3440"}
	ColorFg        lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#24292E", Dark: "#ECEFF4"}
	ColorMuted     lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#6A737D", Dark: "#4C566A"}
)
