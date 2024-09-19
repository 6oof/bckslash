package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectStatusModel struct {
	Viewport viewport.Model
	content  string
	uuid     string
}

func MakeProjectStatusModel(uuid string) ProjectStatusModel {
	return ProjectStatusModel{
		uuid: uuid,
	}
}

func (m ProjectStatusModel) Init() tea.Cmd {
	return nil
}

func (m ProjectStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return MakeProjectModel(), commands.FetchProject(m.uuid)
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.Viewport.Width = constants.BodyWidth()
		m.Viewport.Height = constants.BodyHeight()
		return m, nil

	case commands.ProgramErrMsg:
		return GoError(msg.Err, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.uuid)
		})

	case commands.ExecFinishedMsg:
		m.content = constants.ViewportContent(msg.Content)
		m.Viewport = viewport.New(constants.BodyWidth(), constants.BodyHeight())
		m.Viewport.SetContent(m.content) // Update viewport with the new command output

	case commands.ExecStartMsg:
		return m, commands.ShowProjectStatus(m.uuid)
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m ProjectStatusModel) View() string {

	return constants.Layout("Help", "↑/↓: Scroll • q: Quit", m.Viewport.View())
}
