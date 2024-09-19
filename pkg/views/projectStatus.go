package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectStatusModel struct {
	Err      error
	Viewport viewport.Model
	uuid     string
	Content  string
}

func MakeProjectStatusModel(uuid string) ProjectStatusModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())
	return ProjectStatusModel{
		Err:      nil,
		Viewport: vp,
		Content:  "",
		uuid:     uuid,
	}
}

func (m ProjectStatusModel) Init() tea.Cmd {
	return commands.OpenHelpMd()
}

func (m ProjectStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return MakeProjectModel(), commands.FetchProject(m.uuid)
		case "up", "k":
			m.Viewport.LineUp(1) // Move up
			return m, nil
		case "down", "j":
			m.Viewport.LineDown(1) // Move down
			return m, nil
		case "pageup":
			m.Viewport.LineUp(m.Viewport.Height / 2) // Move up by half the viewport height
			return m, nil
		case "pagedown":
			m.Viewport.LineDown(m.Viewport.Height / 2) // Move down by half the viewport height
			return m, nil
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.Viewport.Width = constants.BodyWidth()
		m.Viewport.Height = constants.BodyHeight()
		m.Viewport.SetContent(m.Content) // Update content on resize

	case commands.ProgramErrMsg:
		return GoError(msg.Err, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.uuid)
		})

	case commands.ExecFinishedMsg:
		m.Viewport.SetContent(msg.Content)
		return m, nil

	case commands.ExecStartMsg:
		return m, commands.ShowProjectStatus(m.uuid)
	}

	return m, nil
}

func (m ProjectStatusModel) View() string {
	if m.Err != nil {
		return constants.Layout("Help", "q: Quit", "Error: "+m.Err.Error()+"\n")
	}
	return constants.Layout("Help", "↑/↓: Scroll • q: Quit", m.Viewport.View())
}
