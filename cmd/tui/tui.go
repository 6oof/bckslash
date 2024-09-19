package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/6oof/bckslash/pkg/helpers"
	"github.com/6oof/bckslash/pkg/views"

	tea "github.com/charmbracelet/bubbletea"
)

// StartTea is the entry point for the UI. Initializes the model.
func StartTea() error {
	// Initialize logging
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	err := helpers.OpenDb()
	if err != nil {
		panic("Falid opening bcks database: " + err.Error())
	}
	defer helpers.CloseDb()

	// Start the TUI application
	p := tea.NewProgram(views.InitHomeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	return nil
}
