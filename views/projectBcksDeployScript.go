package views

import (
	"lg/views/commands"
	"lg/views/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type ProjectBcksDelpoyModel struct {
	Err      error
	Viewport viewport.Model
	uuid     string
	Content  string
}

func MakeProjectBcksDelpoyModel(uuid string) ProjectBcksDelpoyModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())
	return ProjectBcksDelpoyModel{
		Err:      nil,
		Viewport: vp,
		Content:  "",
		uuid:     uuid,
	}
}

func (m ProjectBcksDelpoyModel) Init() tea.Cmd {
	return commands.OpenHelpMd()
}

func (m ProjectBcksDelpoyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.Err = msg.Err

	case commands.ExecFinishedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}
		m.Content = m.renderMarkdown(msg.Content)
		m.Viewport.SetContent(m.Content)
		return m, nil
	case commands.ExecStartMsg:
		return m, commands.OpenProjectBcksDeployScript(m.uuid)
	}

	return m, nil
}

func (m ProjectBcksDelpoyModel) View() string {
	if m.Err != nil {
		return constants.Layout("Help", "q: Quit", "Error: "+m.Err.Error()+"\n")
	}
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
