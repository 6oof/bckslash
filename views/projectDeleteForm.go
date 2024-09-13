package views

import (
	"lg/views/commands"
	"lg/views/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ProjectDeleteModel struct {
	form        *huh.Form // huh.Form is just a tea.Model
	confirm     bool
	projectUuid string
	Err         error
	loading     bool
}

func MakePeojectDeleteModel(uuid string) (ProjectDeleteModel, error) {
	// Load settings

	fm := ProjectDeleteModel{
		confirm:     false,
		Err:         nil,
		projectUuid: uuid,
	}

	// Create form with editor options
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
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			// Return to home on escape
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

		// Return to home after saving
		return GoToProjects()
	}

	return m, cmd
}

func (m ProjectDeleteModel) View() string {
	if m.loading {
		return constants.Layout("Editor Selection", "q: Return home", "Loading...")
	}

	if m.Err != nil {
		return constants.Layout("Editor Selection", "q: Return home", "Error: "+m.Err.Error()+"\n")
	}

	return constants.Layout("Editor Selection", "q: back", constants.PadBodyContent.Render(m.form.View()))
}
