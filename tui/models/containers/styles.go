package containers

import (
	"github.com/charmbracelet/lipgloss"
	style "github.com/kdruelle/gmd/tui/styles"
)

var (
	colNameStyle    = lipgloss.NewStyle().Width(40)
	colStateStyle   = lipgloss.NewStyle().Width(20)
	colImageStyle   = lipgloss.NewStyle().Width(70)
	colAddressStyle = lipgloss.NewStyle().Width(40)
)

var (
	UpdateUnavailable   = style.Inactive().Render("-")
	UpToDateFlag        = style.Success().Render("✓")
	UpdateAvailableFlag = style.Danger().Render("⚠")
)

var (
	ContainerRuningState     = style.Success().Render("running")
	ContainerExitedState     = style.Danger().Render("exited")
	ContainerCreatedState    = style.Inactive().Render("created")
	ContainerPausedState     = style.Inactive().Render("paused")
	ContainerRestartingState = style.Warning().Render("restarting")
)
