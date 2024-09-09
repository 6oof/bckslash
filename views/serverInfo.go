package views

import (
	"lg/views/commands"
	"lg/views/constants"
	"lg/views/layout"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerInfoModel struct {
	Err     error
	Content string
}

func InitServerInfoModel() ServerInfoModel {
	return ServerInfoModel{}
}

func (m ServerInfoModel) Init() tea.Cmd {
	return nil
}

func (m ServerInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
	case commands.ExecStartMsg:
		return m, commands.ShowNeofetch()
	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		} else {
			if msg.Content != "" {
				m.Content = msg.Content
			}
		}
		return m, nil

	}
	return m, nil
}

func (m ServerInfoModel) View() string {
	if m.Err != nil {
		return layout.Layout("Server Info", "q: Return home", "Error: "+m.Err.Error()+"\n")
	}

	return layout.Layout("Server Info", "q: Return home", m.Content)

}
