package views

import (
	"errors"
	"regexp"

	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ProjectAddModel struct {
	form    *huh.Form // huh.Form is just a tea.Model
	Err     error
	loading bool
}

func MakePeojectAddModel() ProjectAddModel {
	// Load settings

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
	).WithTheme(constants.HuhBsTheme())

	return fm
}

func (m ProjectAddModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ProjectAddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			// Return to home on escape
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
		return m, commands.AddProjectCommand(m.form.GetString("name"), "project", m.form.GetString("repo"), m.form.GetString("branch"))
	}

	return m, cmd
}

func (m ProjectAddModel) View() string {
	if m.loading {
		return constants.Layout("Editor Selection", "q: Return home", "Loading...")
	}

	if m.Err != nil {
		return constants.Layout("Editor Selection", "q: Return home", constants.PadBodyContent.Render("Error: "+m.Err.Error()+"\n"))
	}

	return constants.Layout("Editor Selection", "q: back", constants.PadBodyContent.Render(m.form.View()))
}
