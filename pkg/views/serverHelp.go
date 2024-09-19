package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type ServerHelpModel struct {
	Err      error
	Viewport viewport.Model
	Content  string
}

func MakeServerHelpModel() ServerHelpModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())
	return ServerHelpModel{
		Err:      nil,
		Viewport: vp,
		Content:  "",
	}
}

func (m ServerHelpModel) Init() tea.Cmd {
	return commands.OpenHelpMd()
}

func (m ServerHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		case "up", "k":
			m.Viewport.LineUp(3) // Move up
			return m, nil
		case "down", "j":
			m.Viewport.LineDown(3) // Move down
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
		return GoError(msg.Err, GoHome)

	case commands.ExecFinishedMsg:
		m.Content = m.renderMarkdown(msg.Content)
		m.Viewport.SetContent(m.Content)
		return m, nil

	case commands.ExecStartMsg:
		return m, commands.OpenHelpMd()
	}

	return m, nil
}

func (m ServerHelpModel) View() string {
	if m.Err != nil {
		return constants.Layout("Help", "q: Quit", "Error: "+m.Err.Error()+"\n")
	}
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
