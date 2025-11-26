package containers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/kdruelle/gmd/docker/types"
	style "github.com/kdruelle/gmd/tui/styles"
)

type ContainerItem struct {
	id           string
	name         string
	state        container.ContainerState
	actionState  container.ContainerState
	update       *bool
	content      string
	statsContent string

	show bool
}

func NewContainerItem(dc types.Container) ContainerItem {
	c := ContainerItem{
		id:    dc.ID,
		name:  dc.Name,
		state: dc.State.Status,
	}

	// if dc.Stats.ID != "" {

	// 	cpuPct := docker.CPUPercent(dc.Stats)
	// 	// 	//memPct := docker.MemoryPercent(c.Stats)
	// 	usedMem := dc.Stats.MemoryStats.Usage
	// 	limitMem := dc.Stats.MemoryStats.Limit

	// 	cpuBar := htopCpuBar(cpuPct, 20) // 20 = largeur de la barre
	// 	memBar := htopMemBar(usedMem, limitMem, 20)

	// 	statsContent = lipgloss.JoinHorizontal(
	// 		lipgloss.Left,
	// 		cpuBar, "   ", memBar,
	// 	)
	// }

	return c
}

func (c *ContainerItem) RenderContent() {

	title := style.Title.Render(c.Name())
	shortID := style.Subtitle.Render(c.ShortID())

	statsContent := "CPU[ -- ]   RAM[ -- ]"
	c.statsContent = col3Style.Render(statsContent)

	col1 := lipgloss.JoinVertical(lipgloss.Left, title, shortID)
	col2 := lipgloss.JoinHorizontal(lipgloss.Center, c.UpdateFlag(), " ", c.Status())

	col1 = col1Style.Render(col1)
	col2 = col2Style.Render(col2)

	c.content = lipgloss.JoinHorizontal(lipgloss.Center, col1, " ", col2)
}

// func (c *ContainerItem) RenderStats(stats container.StatsResponse) {
// 	if stats.ID == "" {
// 		c.statsContent = "CPU[ -- ]   RAM[ -- ]"
// 	}

// 	cpuPct := docker.CPUPercent(stats)
// 	//memPct := docker.MemoryPercent(c.Stats)
// 	usedMem := stats.MemoryStats.Usage
// 	limitMem := stats.MemoryStats.Limit

// 	cpuBar := htopCpuBar(cpuPct, 20) // 20 = largeur de la barre
// 	memBar := htopMemBar(usedMem, limitMem, 20)

// 	c.statsContent = lipgloss.JoinHorizontal(
// 		lipgloss.Left,
// 		cpuBar, "   ", memBar,
// 	)
// }

func (c ContainerItem) Name() string {
	title := strings.TrimPrefix(c.name, "/")
	return title
}

func (c ContainerItem) ShortID() string {
	shortID := c.id
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}
	return shortID
}

func (c ContainerItem) Status() string {

	if c.actionState != "" {
		switch c.actionState {
		case container.StateRestarting:
			return ContainerRestartingState
		}
		return c.actionState
	}

	switch c.state {
	case container.StateRunning:
		return ContainerRuningState
	case container.StateExited:
		return ContainerExitedState
	case container.StateCreated:
		return ContainerCreatedState
	case container.StatePaused:
		return ContainerPausedState
	case container.StateRestarting:
		return ContainerRestartingState
	default:
		return c.state
	}
}

func (c ContainerItem) UpdateFlag() string {
	if c.update == nil {
		return UpdateUnavailable
	}
	if *c.update {
		return UpdateAvailableFlag
	} else {
		return UpToDateFlag
	}
}

// func (c ContainerItem) Description() string {

// 	shortID := c.ID
// 	if len(shortID) > 12 {
// 		shortID = shortID[:12]
// 	}
// 	shortID = style.Subtitle.Render(shortID)

// 	state := c.State.Status

// 	switch state {
// 	case container.StateRunning:
// 		state = style.ActiveItem.Render("running")
// 	case container.StateExited:
// 		state = style.DangerItem.Render("exited")
// 	case container.StateCreated:
// 		state = style.InactiveItem.Render("created")
// 	case container.StatePaused:
// 		state = style.InactiveItem.Render("created")
// 	case "restarting":
// 		state = style.WarningItem.Render("restarting")
// 	default:
// 		return fmt.Sprintf("%s - %s", c.ID, state)
// 	}

// 	// flag update
// 	update := " "
// 	if c.Update {
// 		update = " " + style.DangerItem.Render("⚠")
// 	}

// 	// return fmt.Sprintf("%-16s           %2s%s", style.Subtitle.Render(c.ID), update, state)
// 	// on retourne une description multi-ligne
// 	return lipgloss.JoinHorizontal(
// 		lipgloss.Top,
// 		shortID,
// 		"   ",
// 		update,
// 		state,
// 	)
// }

func (c ContainerItem) FilterValue() string { return c.Name() }

// func (c ContainerItem) StatsView() string {

// 	if c.dc.Stats.ID == "" {
// 		return "CPU --%  RAM --"
// 	}

// 	cpu := fmt.Sprintf("CPU %2.0f%%", docker.CPUPercent(c.dc.Stats))

// 	ram := fmt.Sprintf("RAM %s (%2.0f%%)",
// 		humanize.Bytes(uint64(c.dc.Stats.MemoryStats.Usage)),
// 		docker.MemoryPercent(c.dc.Stats),
// 	)

// 	return lipgloss.JoinHorizontal(lipgloss.Left, cpu, "   ", ram)
// }

// func (c ContainerItem) StatsView() string {
// 	if c.Stats.ID == "" {
// 		return "CPU[ -- ]   RAM[ -- ]"
// 	}

// 	cpuPct := docker.CPUPercent(c.Stats)
// 	//memPct := docker.MemoryPercent(c.Stats)
// 	usedMem := c.Stats.MemoryStats.Usage
// 	limitMem := c.Stats.MemoryStats.Limit

// 	cpuBar := htopCpuBar(cpuPct, 20) // 20 = largeur de la barre
// 	memBar := htopMemBar(usedMem, limitMem, 20)

// 	return lipgloss.JoinHorizontal(
// 		lipgloss.Left,
// 		cpuBar, "   ", memBar,
// 	)
// }

func htopCpuBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}

	ratio := pct / 100
	filled := int(ratio * float64(width))

	var color lipgloss.AdaptiveColor
	switch {
	case ratio < 0.50:
		color = style.ColorSuccess
	case ratio < 0.80:
		color = style.ColorWarning
	default:
		color = style.ColorDanger
	}

	style := lipgloss.NewStyle().Foreground(color)

	return fmt.Sprintf(
		"CPU[%s%s %3.0f%%]",
		style.Render(strings.Repeat("|", filled)),
		strings.Repeat(" ", width-filled),
		pct,
	)
}

func htopMemBar(used, total uint64, width int) string {
	title := style.ActiveItem.Render("Mem")
	if total == 0 {
		return "Mem[ no limit ]"
	}

	usedGB := float64(used) / (1024 * 1024 * 1024)
	totalGB := float64(total) / (1024 * 1024 * 1024)
	pct := usedGB / totalGB

	// Nombre de barres pleines
	filled := int(float64(width) * pct)
	if filled > width {
		filled = width
	}

	// Couleur selon seuil
	var color lipgloss.AdaptiveColor
	switch {
	case pct < 0.50:
		color = style.ColorSuccess
	case pct < 0.80:
		color = style.ColorWarning
	default:
		color = style.ColorDanger
	}

	style := lipgloss.NewStyle().Foreground(color)

	bars := style.Render(strings.Repeat("|", filled))
	padding := strings.Repeat(" ", width-filled)

	return fmt.Sprintf(
		"%s%s%s%s %.1fG/%.1fG%s",
		title,
		lipgloss.NewStyle().Width(1).Bold(true).Render("["),
		bars,
		padding,
		usedGB,
		totalGB,
		lipgloss.NewStyle().Width(1).Bold(true).Render("]"),
	)
}
