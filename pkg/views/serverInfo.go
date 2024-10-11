package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerInfoModel struct {
	Content string
	Loading bool
	Spinner spinner.Model
}

func MakeServerInfoModel() ServerInfoModel {
	// Initialize the spinner with a default style
	s := spinner.New()
	s.Spinner = spinner.Dot

	return ServerInfoModel{
		Spinner: s,
		Loading: true,
	}
}

func (m ServerInfoModel) Init() tea.Cmd {
	// Start the spinner when the model is initialized
	return nil
}

func (m ServerInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.ModalKeymap.Back):
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ExecStartMsg:
		m.Loading = true
		return m, tea.Batch(m.Spinner.Tick, commands.ShowNeofetch())

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

	case commands.ExecFinishedMsg:
		m.Loading = false
		if msg.Content != "" {
			m.Content = msg.Content
		}
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

func (m ServerInfoModel) View() string {
	if m.Loading {
		// Show the spinner while loading
		return constants.Layout("Server Info", constants.ModalHelpString, constants.PadBodyContent.Render(m.Spinner.View()+" Loading..."))
	}

	return constants.Layout("Server Info", constants.ModalHelpString, lipgloss.PlaceHorizontal(constants.BodyWidth(), lipgloss.Left, constants.PadBodyContent.Render(m.Content)))
}
