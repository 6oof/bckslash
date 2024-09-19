package views

import (
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
	deploy
	viewCompose
	viewStatus
	viewDeployScript
	editEnv
)

type ProjectModel struct {
	Err          error
	project      helpers.Project
	shortGitData string
	deployLog    string
	Menu         list.Model
	Loading      bool
	dataLoading  bool
	Spinner      spinner.Model
}

func MakeProjectModel() ProjectModel {
	// Initialize the spinner with a default style
	s := spinner.New()
	s.Spinner = spinner.Dot

	menuItems := []list.Item{
		item{title: "Deploy", desc: "Trigger deployment", navigation: deploy},
		item{title: "Project status", desc: "View docker-compose ps and git status out", navigation: viewStatus},
		item{title: "View compose", desc: "View the bckslash-compose.yaml file", navigation: viewCompose},
		item{title: "View deploy script", desc: "View the bckslash-deploy.sh file", navigation: viewDeployScript},
		item{title: "Enviroment", desc: "Edit .env file", navigation: editEnv},
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
	// Start the spinner when the model is initialized
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return GoHome()
		}
		switch {
		case key.Matches(msg, constants.Keymap.Enter):
			switch m.Menu.SelectedItem().(item).navigation {
			case deleteProject:
				pdm, _ := MakePeojectDeleteModel(m.project.UUID)
				pdm.Init()
				return pdm, nil
			case editEnv:
				return MakeProjectEnvModel(m.project.UUID).Update(commands.ExecStartMsg{})
			case viewCompose:
				return MakeProjectBcksComposeModel(m.project.UUID).Update(commands.ExecStartMsg{})
			case viewDeployScript:
				return MakeProjectBcksDelpoyModel(m.project.UUID).Update(commands.ExecStartMsg{})
			case viewStatus:
				return MakeProjectStatusModel(m.project.UUID).Update(commands.ExecStartMsg{})
			case deploy:
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
		m.deployLog = msg.Content
		return m, nil

	case commands.ProgramErrMsg:
		return GoError(msg.Err, func() (tea.Model, tea.Cmd) {
			return MakeProjectModel(), commands.FetchProject(m.project.UUID)
		})

	case spinner.TickMsg:
		if m.Loading {
			// Update the spinner if loading is true
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
	if m.deployLog != "" {
		return constants.Layout("Server Info", "q: Return home", constants.PadBodyContent.Render(m.deployLog))
	}
	if m.Loading {
		// Show the spinner while loading
		return constants.Layout("Server Info", "q: Return home", constants.PadBodyContent.Render(m.Spinner.View()+" Loading..."))
	}

	if m.Err != nil {
		return constants.Layout("Server Info", "q: Return home", constants.PadBodyContent.Render("Error: "+m.Err.Error()+"\n"))
	}

	crd := ""
	if m.dataLoading {
		crd = m.Spinner.View() + " Loading..."
	} else {
		crd = "UUID:\n " + m.project.UUID + " \n\nActive git commit:\n" + m.shortGitData
	}

	m.Menu.SetSize(constants.BodyHalfWidth(), constants.BodyHeight())
	return constants.Layout(
		"Server Info", "q: Return home",
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
