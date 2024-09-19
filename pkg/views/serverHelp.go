package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type ServerHelpModel struct {
	Viewport viewport.Model
	Content  string
}

func MakeServerHelpModel() ServerHelpModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())
	return ServerHelpModel{
		Viewport: vp,
		Content:  "",
	}
}

func (m ServerHelpModel) Init() tea.Cmd {
	return nil
}

func (m ServerHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.Viewport.Width = constants.BodyWidth()
		m.Viewport.Height = constants.BodyHeight()
		m.Viewport.SetContent(m.Content) // Update content on resize

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

	case commands.ExecFinishedMsg:
		m.Content = m.renderMarkdown(msg.Content)
		m.Viewport.SetContent(m.Content)
		return m, nil

	case commands.ExecStartMsg:
		return m, commands.OpenHelpMd()
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m ServerHelpModel) View() string {
	return constants.Layout("Help", "↑/↓: Scroll • q: Quit", m.Viewport.View())
}

func (m ServerHelpModel) renderMarkdown(content string) string {
	// Glamour rendering with a dark theme for markdown
	out, err := glamour.Render(content, "tokyo-night")
	if err != nil {
		return "Error rendering markdown"
	}
	return out
}
