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
	mainNav  list.Model
	loading  bool
	quitting bool
}

func InitHomeModel() homeModel {
	itemsLeft := []list.Item{
		item{title: "Projects", desc: "Open project list", navigation: listProjects},
		item{title: "Add a project", desc: "Add a project", navigation: addProject},
		item{title: "Server info", desc: "Server monitoring dashboard", navigation: serverInfo},
		item{title: "Server stats", desc: "Server monitoring dashboard", navigation: serverStats},
		item{title: "Editor settings", desc: "Pick prefered text editor to use when needed", navigation: serverSettings},
		item{title: "Help", desc: "Basic information on bckslash", navigation: help},
	}

	m := homeModel{
		mainNav:  list.New(itemsLeft, constants.CustomDelegate(), 40, 30),
		quitting: false,
	}

	m.mainNav.Styles = constants.ListStyle()

	m.mainNav.Title = "Bckslash actions"
	m.mainNav.SetShowHelp(false)

	return m
}

func GoHome() (tea.Model, tea.Cmd) {
	homeModel := InitHomeModel()
	return homeModel.Update(constants.WinSize)
}

func (m homeModel) Init() tea.Cmd {
	return nil
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			switch m.mainNav.SelectedItem().(item).navigation {
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
	m.mainNav, cmd = m.mainNav.Update(msg)

	return m, cmd
}

func (m homeModel) View() string {
	homeCard := "Welcome to Bckslash!\n\nget help at: github.com/6oof/bckslash"

	if m.quitting {
		return ""
	}

	// Set the width for the lists and render them
	if m.loading {
		return constants.Layout("home", "", "loading...")
	} else {
		m.mainNav.SetSize(constants.BodyHalfWidth(), constants.BodyHeight())
		return constants.Layout(
			"home",
			constants.MainHelpString,
			constants.HalfAndHalfComposition(
				m.mainNav.View(),
				constants.Card(
					homeCard,
					`\`,
					constants.BodyHalfWidth(),
					constants.BodyHeight(),
				),
			),
		)

	}
}
