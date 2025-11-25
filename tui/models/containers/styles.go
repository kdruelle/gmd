package containers

import style "github.com/kdruelle/gmd/tui/styles"

var (
	UpToDateFlag        = style.ActiveItem.Render("✓")
	UpdateAvailableFlag = style.DangerItem.Render("⚠")
)

var (
	ContainerRuningState = style.ActiveItem.Render("running")
	ContainerExitedState = style.DangerItem.Render("exited")
	ContainerCreatedState = style.InactiveItem.Render("created")
	ContainerPausedState = style.InactiveItem.Render("paused")
	ContainerRestartingState = style.WarningItem.Render("restarting")
)