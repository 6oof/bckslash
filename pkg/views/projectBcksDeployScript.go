package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type ProjectBcksDelpoyModel struct {
	Viewport viewport.Model
	uuid     string
}

func MakeProjectBcksDelpoyModel(uuid string) ProjectBcksDelpoyModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())

	return ProjectBcksDelpoyModel{
		Viewport: vp,
		uuid:     uuid,
	}
}

func (m ProjectBcksDelpoyModel) Init() tea.Cmd {
	return nil
}

func (m ProjectBcksDelpoyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case commands.ProgramErrMsg:
		return GoError(msg.Err, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.uuid)
		})

	case commands.ExecFinishedMsg:
		m.Viewport.SetContent(m.renderMarkdown(msg.Content))
		return m, nil

	case commands.ExecStartMsg:
		return m, commands.OpenProjectBcksDeployScript(m.uuid)
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m ProjectBcksDelpoyModel) View() string {

	return constants.Layout("Help", "↑/↓: Scroll • q: Quit", m.Viewport.View())
}

func (m ProjectBcksDelpoyModel) renderMarkdown(content string) string {
	// Glamour rendering with a dark theme for markdown
	cntnt := "```bash\n" + content + "```"

	out, err := glamour.Render(cntnt, "tokyo-night")
	if err != nil {
		return "Error rendering markdown"
	}
	return out
}
