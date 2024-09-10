package views

import (
	"lg/views/commands"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerStatsModel struct {
	Err error
}

func MakeServerStatsModel() ServerStatsModel {
	return ServerStatsModel{Err: nil}
}

func (m ServerStatsModel) Init() tea.Cmd {
	return nil
}

func (m ServerStatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
	case commands.ExecStartMsg:
		return m, commands.OpenBTM()
	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		return GoHome()
	}
	return m, nil
}

func (m ServerStatsModel) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n"
	}
	return "Press 'q' to quit."
}
