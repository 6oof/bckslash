package views

import (
	"github.com/6oof/bckslash/pkg/constants"
	"github.com/charmbracelet/bubbles/key"
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
	vp.SetContent(constants.ViewportContent("Error: " + err.Error()))

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
		switch {
		case key.Matches(msg, constants.ModalKeymap.Back):
			return m.ReturnMc()
		}

	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.Viewport.Width = constants.BodyWidth()
		m.Viewport.Height = constants.BodyHeight()
		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ErrorModel) View() string {
	return constants.Layout("ERROR", constants.ModalHelpString, m.Viewport.View())
}
