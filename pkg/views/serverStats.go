package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerStatsModel struct {
}

func MakeServerStatsModel() ServerStatsModel {
	return ServerStatsModel{}
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
		constants.WinSize = msg

	case commands.ExecStartMsg:
		return m, commands.OpenBTM()

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

	case commands.ExecFinishedMsg:
		return GoHome()

	}
	return m, nil
}

func (m ServerStatsModel) View() string {
	return "Press 'q' to quit."
}
