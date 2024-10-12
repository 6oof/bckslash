package views

import (
	"errors"
	"regexp"

	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ProjectAddModel struct {
	form    *huh.Form
	Err     error
	loading bool
}

func MakePeojectAddModel() ProjectAddModel {
	fm := ProjectAddModel{
		Err: nil,
	}

	// Create form with editor options
	fm.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("Must only contain letters and numbers.").
				Validate(func(str string) error {
					re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
					if !re.MatchString(str) {
						return errors.New("input must only contain letters and numbers, no special characters or spaces")
					}
					return nil
				}),
			huh.NewInput().
				Key("repo").
				Title("Repository link").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("can't be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("branch").
				Title("Branch to clone").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("can't be empty")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Key("servicename").
				Title("Name of the docker-compose service").
				Description("Bckslash merges the traefik labels on the service that's exposed to the internet.").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("can't be empty")
					}

					re := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
					if !re.MatchString(str) {
						return errors.New("invalid service name: must contain only lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen or have consecutive hyphens")
					}

					return nil
				}),
			huh.NewInput().
				Key("domain").
				Title("Domain for the project").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("can't be empty")
					}

					re := regexp.MustCompile(`^([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)$`)
					if !re.MatchString(str) {
						return errors.New("invalid domain format (e.g., example.com)")
					}

					return nil
				}),
		),
	).WithTheme(constants.HuhBsTheme())

	return fm
}

func (m ProjectAddModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ProjectAddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.FormKeymap.Back):
			homeModel := InitHomeModel()
			return homeModel.Update(constants.WinSize)
		}

	case tea.WindowSizeMsg:
		constants.WinSize = msg

	case commands.ProgramErrMsg:
		return GoError(msg.Err, GoHome)

	case commands.ProjectListChangedMsg:
		return GoHome()

	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.loading = true
		m.form.State = huh.StateAborted
		return m, commands.AddProjectCommand(m.form.GetString("name"), "project", m.form.GetString("repo"), m.form.GetString("branch"), m.form.GetString("servicename"), m.form.GetString("domain"))
	}

	return m, cmd
}

func (m ProjectAddModel) View() string {
	if m.loading {
		return constants.Layout("Editor Selection", constants.FormHelpString, "Loading...")
	}

	if m.Err != nil {
		return constants.Layout("Editor Selection", constants.FormHelpString, constants.PadBodyContent.Render("Error: "+m.Err.Error()+"\n"))
	}

	return constants.Layout("Editor Selection", constants.FormHelpString, constants.PadBodyContent.Render(m.form.View()))
}
