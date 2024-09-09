package views

import (
	"fmt"
	"lg/views/constants"
	"lg/views/layout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Model struct {
	form *huh.Form // huh.Form is just a tea.Model
}

func NewTestModel() Model {
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("class").
					Options(huh.NewOptions("Warrior", "Mage", "Rogue")...).
					Title("Choose your class"),

				huh.NewSelect[int]().
					Key("level").
					Options(huh.NewOptions(1, 20, 9999)...).
					Title("Choose your level"),
			),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ...
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			homeModel := InitHomeModel()
			return homeModel.Update(constants.WinSize)
		}

	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m Model) View() string {
	if m.form.State == huh.StateCompleted {
		class := m.form.GetString("class")
		level := m.form.GetString("level")
		return fmt.Sprintf("You selected: %s, Lvl. %d", class, level)
	}

	return layout.Layout("Form", "↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit", m.form.View())

}
