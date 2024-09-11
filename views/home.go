package views

import (
	"lg/views/commands"
	"lg/views/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type navigate int

const (
	serverStats navigate = iota
	listProjects
	addProject
	serverInfo
	serverSettings
	help
	no
)

type item struct {
	title, desc string
	navigation  navigate
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type homeModel struct {
	listLeft list.Model
	loading  bool
	quitting bool
}

func InitHomeModel() homeModel {
	itemsLeft := []list.Item{
		item{title: "List Projects", desc: "Open project list", navigation: listProjects},
		item{title: "Add a project", desc: "Add a project", navigation: addProject},
		item{title: "Server info", desc: "Server monitoring dashboard", navigation: serverInfo},
		item{title: "Server stats", desc: "Server monitoring dashboard", navigation: serverStats},
		item{title: "Edito settings", desc: "Pick prefered text editor to use when needed", navigation: serverSettings},
		item{title: "Help", desc: "Basic information on bckslash", navigation: help},
	}

	m := homeModel{
		listLeft: list.New(itemsLeft, constants.FocusedListDelegate(), 20, 20),
		quitting: false,
	}

	m.listLeft.Styles = constants.ListStyle()

	m.listLeft.Title = "Bckslash actions"
	m.listLeft.SetShowHelp(false)

	return m
}

func GoHome() (tea.Model, tea.Cmd) {
	homeModel := InitHomeModel()
	return homeModel.Update(constants.WinSize)
}

func (m homeModel) Init() tea.Cmd {
	return commands.LoadProjectsCmd()
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			switch m.listLeft.SelectedItem().(item).navigation {
			case addProject:
				apm := MakePeojectAddModel()
				return apm, apm.Init()
			case listProjects:
				return MakeProjectsModel(), commands.LoadProjectsCmd()
			case serverStats:
				return MakeServerStatsModel().Update(commands.ExecStartMsg{})

			case serverInfo:
				return MakeServerInfoModel().Update(commands.ExecStartMsg{})

			case help:
				return MakeServerHelpModel().Update(commands.ExecStartMsg{})

			case serverSettings:
				esm := MakeEditorSelectionModel()
				return esm, esm.Init()

			}
		case key.Matches(msg, constants.Keymap.Back):
			return m, nil

		}
	case commands.ProgramErrMsg:
		panic(msg.Err)

	case tea.WindowSizeMsg:
		constants.WinSize = msg
	}

	var cmd tea.Cmd
	m.listLeft, cmd = m.listLeft.Update(msg)

	return m, cmd
}

func (m homeModel) View() string {

	if m.quitting {
		return ""
	}

	// Set the width for the lists and render them
	if m.loading {
		return constants.Layout("home", "", "loading...")
	} else {
		m.listLeft.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		return constants.Layout("home", constants.MainHelpString, m.listLeft.View())

	}
}
