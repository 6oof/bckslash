package views

import (
	"lg/helpers"
	"lg/views/commands"
	"lg/views/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var projectsState []helpers.Project

type navigate int

const (
	serverStats navigate = iota
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

type project struct {
	title, desc string
	uuid        string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (i project) Title() string       { return i.title }
func (i project) Description() string { return i.desc }
func (i project) FilterValue() string { return i.title }

type homeModel struct {
	listLeft  list.Model
	listRight list.Model
	loading   bool
	quitting  bool
	focus     int // 0 for left list, 1 for right list
}

func InitHomeModel() homeModel {
	itemsLeft := []list.Item{
		item{title: "Add project", desc: "Add a project", navigation: addProject},
		item{title: "Server info", desc: "Server monitoring dashboard", navigation: serverInfo},
		item{title: "Server stats", desc: "Server monitoring dashboard", navigation: serverStats},
		item{title: "Edito settings", desc: "Pick prefered text editor to use when needed", navigation: serverSettings},
		item{title: "Help", desc: "Basic information on bckslash", navigation: help},
	}

	itemsRight := []list.Item{}

	// Populate the new items from the project list
	for _, p := range projectsState {
		itemsRight = append(itemsRight, project{title: p.Title, desc: p.UUID, uuid: p.UUID})
	}

	m := homeModel{
		listLeft:  list.New(itemsLeft, constants.UnfocusedListDelegate(), 20, 20),
		listRight: list.New(itemsRight, constants.UnfocusedListDelegate(), 20, 20),
		focus:     0, // Start with focus on the left list
		quitting:  false,
	}

	m.listLeft.Styles = constants.ListStyle()
	m.listRight.Styles = constants.ListStyle()

	m.listLeft.Title = "Bckslash actions"
	m.listRight.Title = "Projects"
	m.listRight.SetShowHelp(false)
	m.listLeft.SetShowHelp(false)

	return m
}

func (m *homeModel) mapProjects(projects []helpers.Project) {
	itemsRight := []list.Item{}

	// Populate the new items from the project list
	for _, p := range projects {
		itemsRight = append(itemsRight, project{title: p.Title, desc: p.UUID, uuid: p.UUID})
	}

	// Set the new list of items in the right list
	m.listRight.SetItems(itemsRight)
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
		case key.Matches(msg, constants.Keymap.Tab):
			// Switch focus between lists
			m.focus = (m.focus + 1) % 2
			return m, nil
		case key.Matches(msg, constants.Keymap.Enter):
			switch m.focus {
			case 0:
				switch m.listLeft.SelectedItem().(item).navigation {
				case addProject:
					apm := MakePeojectAddModel()
					return apm, apm.Init()
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

			case 1:
			}
		case key.Matches(msg, constants.Keymap.Delete):
			if m.focus == 1 {
				pdm, _ := MakePeojectDeleteModel(m.listRight.SelectedItem().(project).uuid)
				pdm.Init()
				return pdm, nil
			}
		case key.Matches(msg, constants.Keymap.Back):
			return m, nil

		}
	case commands.ProgramErrMsg:
		// setup an err view and handle arr errors with it
		panic(msg.Err)
	case commands.ProjectListChangedMsg:
		m.loading = true
		m.mapProjects(msg.ProjectList)
		projectsState = msg.ProjectList
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		constants.WinSize = msg
	}

	// Update the focused list
	var cmd tea.Cmd
	if m.focus == 0 {
		m.listLeft, cmd = m.listLeft.Update(msg)
	} else {
		m.listRight, cmd = m.listRight.Update(msg)
	}
	return m, cmd
}

func (m homeModel) View() string {

	if m.quitting {
		return ""
	}

	if m.focus == 0 {
		m.listLeft.SetDelegate(constants.FocusedListDelegate())
	} else {
		m.listRight.SetDelegate(constants.FocusedListDelegate())
	}

	// Set the width for the lists and render them
	if m.loading {
		return constants.Layout("home", "", "loading...")
	} else {
		m.listLeft.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		m.listRight.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		return constants.Layout("home", constants.MainHelpString, constants.HalfAndHalfComposition(m.listLeft.View(), m.listRight.View(), constants.BodyHeight()))

	}
}
