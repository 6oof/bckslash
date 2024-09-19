package views

import (
	"path"

	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	tea "github.com/charmbracelet/bubbletea"
)

type ProjectEnvModel struct {
	Err  error
	uuid string
}

func MakeProjectEnvModel(uuid string) ProjectEnvModel {
	return ProjectEnvModel{
		Err:  nil,
		uuid: uuid,
	}
}

func (m ProjectEnvModel) Init() tea.Cmd {
	return nil
}

func (m ProjectEnvModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
	case commands.ExecStartMsg:
		return m, commands.OpenEditor(path.Join("projects", m.uuid, ".env"))
	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		return MakeProjectModel(), commands.FetchProject(m.uuid)
	}
	return m, nil
}

func (m ProjectEnvModel) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n"
	}
	return "Press 'q' to quit."
}
