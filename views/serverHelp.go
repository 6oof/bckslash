package views

import (
	"lg/views/commands"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerHelpModel struct {
	Err error
}

func (m ServerHelpModel) Init() tea.Cmd {
	return nil
}

func InitServerHelpModel() ServerHelpModel {
	return ServerHelpModel{Err: nil}
}

func (m ServerHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
	case commands.ExecStartMsg:
		return m, commands.OpenGlowHelp()
	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		return GoHome()
	}
	return m, nil
}

func (m ServerHelpModel) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n"
	}
	return "Press 'q' to quit."
}
