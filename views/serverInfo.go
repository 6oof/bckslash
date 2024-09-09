package views

import (
	"lg/views/commands"
	"lg/views/constants"
	"lg/views/layout"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ServerInfoModel struct {
	Err     error
	Content string
	Loading bool
	Spinner spinner.Model
}

func InitServerInfoModel() ServerInfoModel {
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
		switch msg.String() {
		case "ctrl+c", "q":
			return GoHome()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ExecStartMsg:
		m.Loading = true
		return m, tea.Batch(m.Spinner.Tick, commands.ShowNeofetch())

	case commands.ExecFinishedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}
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
		return layout.Layout("Server Info", "q: Return home", m.Spinner.View()+" Loading...")
	}

	if m.Err != nil {
		return layout.Layout("Server Info", "q: Return home", "Error: "+m.Err.Error()+"\n")
	}

	return layout.Layout("Server Info", "q: Return home", m.Content)
}
