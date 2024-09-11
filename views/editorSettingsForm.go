package views

import (
	"fmt"
	"lg/helpers"
	"lg/views/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type EditorFormModel struct {
	form    *huh.Form // huh.Form is just a tea.Model
	confirm bool
	Err     error
}

func MakeEditorSelectionModel() EditorFormModel {
	// Load settings
	settings, err := helpers.GetSettings()
	if err != nil {
		return EditorFormModel{Err: fmt.Errorf("failed to load settings: %w", err)}
	}

	fm := EditorFormModel{
		confirm: false,
		Err:     nil,
	}

	// Create form with editor options
	fm.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("editor").
				Description("Some action within bckslash will require you to edit docker-compose, json, and other files. This sets the editor you use for these actions.").
				Options(huh.NewOptions("vim", "nano", "micro")...). // Options for the editor
				Title("Choose your prefered editor:").Value(&settings.EditorCommand),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&fm.confirm),
		),
	)

	return fm
}

func (m EditorFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m EditorFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		m.confirm = f.GetBool("confirm")
	}

	if m.form.State == huh.StateCompleted {
		if m.confirm {
			// Save the selected editor to settings
			selectedEditor := m.form.GetString("editor")

			// Load settings, modify, and save
			settings, err := helpers.GetSettings()
			if err != nil {
				m.Err = err
				return m, nil
			}
			settings.EditorCommand = selectedEditor

			// Save updated settings
			if err := helpers.SaveSettings(settings); err != nil {
				m.Err = err
				return m, nil
			}

		}

		// Return to home after saving
		homeModel := InitHomeModel()
		return homeModel.Update(constants.WinSize)
	}

	return m, cmd
}

func (m EditorFormModel) View() string {
	if m.Err != nil {
		return constants.Layout("Editor Selection", "q: Return home", constants.PadBodyContent.Render("Error: "+m.Err.Error()+"\n"))
	}

	return constants.Layout("Editor Selection", "q: back", constants.PadBodyContent.Render(m.form.View()))
}
