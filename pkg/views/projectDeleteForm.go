package views

import (
	"github.com/6oof/bckslash/pkg/commands"
	"github.com/6oof/bckslash/pkg/constants"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ProjectDeleteModel struct {
	form        *huh.Form
	confirm     bool
	projectUuid string
	Err         error
	loading     bool
}

func MakePeojectDeleteModel(uuid string) (ProjectDeleteModel, error) {

	fm := ProjectDeleteModel{
		confirm:     false,
		Err:         nil,
		projectUuid: uuid,
	}

	fm.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Description("this action delets the project from the list but does not remove any directories or stop docker processes").
				Title("Are you sure you want to delte project " + fm.projectUuid).
				Affirmative("Yes!").
				Negative("No.").
				Value(&fm.confirm),
		),
	).WithTheme(constants.HuhBsTheme())

	return fm, nil
}

func (m ProjectDeleteModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ProjectDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.FormKeymap.Back):
			return GoToProjects()
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		m.confirm = f.GetBool("confirm")
	}

	if m.form.State == huh.StateCompleted {

		if m.confirm {
			m.loading = true
			return MakeProjectsModel(), commands.DeleteProjectCommand(m.projectUuid)
		}

		return GoToProjects()
	}

	return m, cmd
}

func (m ProjectDeleteModel) View() string {
	if m.loading {
		return constants.Layout("Delete Project", constants.FormHelpString, "Loading...")
	}

	if m.Err != nil {
		return constants.Layout("Delete Project", constants.FormHelpString, "Error: "+m.Err.Error()+"\n")
	}

	return constants.Layout("Delete Project", constants.FormHelpString, constants.PadBodyContent.Render(m.form.View()))
}
