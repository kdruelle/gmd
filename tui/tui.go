package tui

import (
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func Start(debugFile string) (err error) {

	if debugFile != "" {
		f, err := tea.LogToFile(debugFile, "debug")
		if err != nil {
			log.Fatalf("unable to open debug file %s : %s", debugFile, err)
		}
		log.SetOutput(f)
		defer func() { err = f.Close() }()
	} else {
		log.SetOutput(io.Discard)
	}

	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Erreur au lancement du TUI : %v", err)
		return err
	}
	return
}
