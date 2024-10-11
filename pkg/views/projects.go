package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"
	"github.com/6oof/bckslash/pkg/helpers"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var projectsState []helpers.Project

type project struct {
	title, desc string
	uuid        string
}

func (i project) Title() string       { return i.title }
func (i project) Description() string { return i.desc }
func (i project) FilterValue() string { return i.title }

type projectsModel struct {
	projectList list.Model
	loading     bool
	quitting    bool
}

func MakeProjectsModel() projectsModel {
	itemsRight := []list.Item{}

	// Populate the new items from the project list
	for _, p := range projectsState {
		itemsRight = append(itemsRight, project{title: p.Title, desc: p.UUID, uuid: p.UUID})
	}

	m := projectsModel{
		projectList: list.New(itemsRight, constants.CustomDelegate(), 20, 20),
		quitting:    false,
		loading:     true,
	}

	m.projectList.Styles = constants.ListStyle()

	m.projectList.Title = "Projects"
	m.projectList.SetShowHelp(false)

	return m
}

func GoToProjects() (tea.Model, tea.Cmd) {
	homeModel := MakeProjectsModel()
	return homeModel, commands.LoadProjectsCmd()
}

func (m *projectsModel) mapProjects(projects []helpers.Project) {
	itemsRight := []list.Item{}

	// Populate the new items from the project list
	for _, p := range projects {
		itemsRight = append(itemsRight, project{title: p.Title, desc: p.UUID, uuid: p.UUID})
	}

	// Set the new list of items in the right list
	m.projectList.SetItems(itemsRight)
}

func (m projectsModel) Init() tea.Cmd {
	return commands.LoadProjectsCmd()
}

func (m projectsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.HomeKeymap.Quit):
			return GoHome()

		case key.Matches(msg, constants.HomeKeymap.Back):
			return GoHome()

		case key.Matches(msg, constants.HomeKeymap.Enter):
			uuid := m.projectList.SelectedItem().(project).uuid
			return MakeProjectModel(), commands.FetchProject(uuid)
		}

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

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
	m.projectList, cmd = m.projectList.Update(msg)
	return m, cmd
}

func (m projectsModel) View() string {

	if m.quitting {
		return ""
	}

	// Set the width for the lists and render them
	if m.loading {
		return constants.Layout("home", "", "loading...")
	} else {
		m.projectList.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		return constants.Layout("home", constants.HomeHelpString, m.projectList.View())

	}
}
