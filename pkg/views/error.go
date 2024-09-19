package views

import (
	"github.com/6oof/bckslash/pkg/constants"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type errReturnModel func() (tea.Model, tea.Cmd)

type ErrorModel struct {
	ReturnMc func() (tea.Model, tea.Cmd)
	Viewport viewport.Model
}

func makeErrorModel(err error, returnMc errReturnModel) ErrorModel {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())

	vp.Width = constants.BodyWidth()
	vp.Height = constants.BodyHeight()
	vp.Style = vp.Style.Padding(2, 0)

	vp.SetContent(err.Error())

	return ErrorModel{
		ReturnMc: returnMc,
		Viewport: vp,
	}
}

func GoError(err error, returnFun errReturnModel) (tea.Model, tea.Cmd) {
	homeModel := makeErrorModel(err, returnFun)
	return homeModel, nil
}

func (m ErrorModel) Init() tea.Cmd {
	return nil
}

func (m ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.ReturnMc()
		case "up", "k":
			m.Viewport.LineUp(1) // Move up
			return m, nil
		case "down", "j":
			m.Viewport.LineDown(1) // Move down
			return m, nil
		case "pageup":
			m.Viewport.LineUp(m.Viewport.Height / 2) // Move up by half the viewport height
			return m, nil
		case "pagedown":
			m.Viewport.LineDown(m.Viewport.Height / 2) // Move down by half the viewport height
			return m, nil
		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.Viewport.Width = constants.BodyWidth()
		m.Viewport.Height = constants.BodyHeight()
		return m, nil
	}

	return m, nil
}

func (m ErrorModel) View() string {
	return constants.Layout("ERROR", "q: Quit", m.Viewport.View())
}
