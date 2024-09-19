package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerSettingsModel struct {
	Err error
}

func MakeServerSettingsModel() ServerSettingsModel {
	return ServerSettingsModel{Err: nil}
}

func (m ServerSettingsModel) Init() tea.Cmd {
	return nil
}

func (m ServerSettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
	case commands.ExecStartMsg:
		return m, commands.OpenEditor("bckslash_settings.json")
	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		return GoHome()
	}
	return m, nil
}

func (m ServerSettingsModel) View() string {
	if m.Err != nil {
		return constants.Layout("Server Info", "q: Return home", "Error: "+m.Err.Error()+"\n")
	}
	return "Press 'q' to quit."
}
