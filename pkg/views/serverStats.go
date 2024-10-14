package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/key"
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
		switch {
		case key.Matches(msg, constants.FormKeymap.Back):
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ExecStartMsg:
		return m, commands.OpenHtop()

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

	case commands.ExecFinishedMsg:
		return GoHome()

	}
	return m, nil
}

func (m ServerStatsModel) View() string {
	return constants.FormHelpString
}
