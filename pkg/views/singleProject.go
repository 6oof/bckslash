package views

import (
	"path"

	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"
	"github.com/6oof/bckslash/pkg/helpers"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	deleteProject navigate = iota
	deployProject
	viewStatusProject
	viewDeployScriptProject
	editEnvProject
	editProxyProject
	executeCommandProject
)

type ProjectModel struct {
	project      helpers.Project
	shortGitData string
	Menu         list.Model
	Loading      bool
	dataLoading  bool
	Spinner      spinner.Model
}

func MakeProjectModel() ProjectModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	menuItems := []list.Item{
		item{title: "Deploy", desc: "Trigger deployment", navigation: deployProject},
		item{title: "Project status", desc: "View docker-compose ps and git status out", navigation: viewStatusProject},
		item{title: "Execute commands", desc: "Opens a shell in project directory", navigation: executeCommandProject},
		item{title: "View actions", desc: "View the bckslash-actions.sh file", navigation: viewDeployScriptProject},
		item{title: "Domain", desc: "Edit reverse proxy labels", navigation: editProxyProject},
		item{title: "Enviroment", desc: "Edit .env file", navigation: editEnvProject},
		item{title: "Delete a project", desc: "You'll be asked to confirm", navigation: deleteProject},
	}

	m := ProjectModel{
		Spinner:     s,
		Loading:     true,
		dataLoading: true,
		Menu:        list.New(menuItems, constants.CustomDelegate(), constants.BodyHalfWidth(), constants.BodyHeight()),
	}

	m.Menu.Styles = constants.ListStyle()

	m.Menu.Title = "Project actions"
	m.Menu.SetShowHelp(false)

	return m
}

func (m ProjectModel) Init() tea.Cmd {
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.HomeKeymap.Back):
			return GoHome()

		case key.Matches(msg, constants.HomeKeymap.Enter):
			switch m.Menu.SelectedItem().(item).navigation {
			case deleteProject:
				pdm, _ := MakePeojectDeleteModel(m.project.UUID)
				pdm.Init()
				return pdm, nil
			case editEnvProject:
				return m, commands.OpenEditor(path.Join(constants.ProjectsDir, m.project.UUID, ".env"))
			case editProxyProject:
				return m, commands.OpenEditor(path.Join(constants.ProjectsDir, m.project.UUID, ".bckslash", "bckslash-traefik-compose.yaml"))
			case executeCommandProject:
				return m, commands.ShellInProject(m.project.UUID)
			case viewDeployScriptProject:
				m.Loading = true
				return m, commands.OpenProjectBcksDeployScript(m.project.UUID)
			case viewStatusProject:
				m.Loading = true
				return m, commands.ShowProjectStatus(m.project.UUID)
			case deployProject:
				return m, commands.TriggerDeploy(m.project.UUID)
			}
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ProjectFoundMsg:
		m.Loading = false
		m.project = msg.Project
		return m, commands.FetchProjectData(m.project.UUID)

	case commands.ProjectViewData:
		m.shortGitData = msg.GitLog
		m.dataLoading = false
		return m, nil

	case commands.ExecFinishedMsg:
		return GoSuccess(msg.Content, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.project.UUID)
		})

	case commands.ProgramErrMsg:
		return GoError(msg.Err, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.project.UUID)
		})

	case spinner.TickMsg:
		if m.Loading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.Menu, cmd = m.Menu.Update(msg)

	return m, cmd
}

func (m ProjectModel) View() string {
	if m.Loading {
		return constants.Layout("Project", constants.HomeAltHelpString, constants.PadBodyContent.Render(m.Spinner.View()+" Loading..."))
	}

	crd := ""
	if m.dataLoading {
		crd = m.Spinner.View() + " Loading..."
	} else {
		crd = "Title:\n " + m.project.Title + " \n\nUUID:\n " + m.project.UUID + " \n\nActive git commit:\n" + m.shortGitData
	}

	m.Menu.SetSize(constants.BodyHalfWidth(), constants.BodyHeight())
	return constants.Layout(
		"Project", constants.HomeAltHelpString,
		lipgloss.PlaceHorizontal(
			constants.BodyWidth(),
			lipgloss.Left,
			constants.HalfAndHalfComposition(
				m.Menu.View(),
				constants.Card(
					crd,
					`.`,
					constants.BodyHalfWidth(),
					constants.BodyHeight(),
				),
			),
		),
	)
}
