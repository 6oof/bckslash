package views

import (
	"lg/helpers"
	"lg/views/commands"
	"lg/views/constants"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectModel struct {
	Err     error
	project helpers.Project
	Loading bool
	Spinner spinner.Model
}

func (m *ProjectModel) inti() {}

func MakeProjectModel() ProjectModel {
	// Initialize the spinner with a default style
	s := spinner.New()
	s.Spinner = spinner.Dot

	return ProjectModel{
		Spinner: s,
		Loading: true,
	}
}

func (m ProjectModel) Init() tea.Cmd {
	// Start the spinner when the model is initialized
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ProjectFoundMsg:
		m.Loading = false
		m.project = msg.Project
		return m, nil

	case spinner.TickMsg:
		if m.Loading {
			// Update the spinner if loading is true
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m ProjectModel) View() string {
	if m.Loading {
		// Show the spinner while loading
		return constants.Layout("Server Info", "q: Return home", constants.PadBodyContent.Render(m.Spinner.View()+" Loading..."))
	}

	if m.Err != nil {
		return constants.Layout("Server Info", "q: Return home", constants.PadBodyContent.Render("Error: "+m.Err.Error()+"\n"))
	}

	return constants.Layout("Server Info", "q: Return home", lipgloss.PlaceHorizontal(constants.BodyWidth(), lipgloss.Left, constants.PadBodyContent.Render(m.project.UUID)))
}
