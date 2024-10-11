package views

import (
	"github.com/6oof/bckslash/pkg/constants"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type successReturnModel func() (tea.Model, tea.Cmd)

type SuccessModle struct {
	ReturnMc func() (tea.Model, tea.Cmd)
	Viewport viewport.Model
}

func makeSuccessModle(content string, returnMc errReturnModel) SuccessModle {
	vp := viewport.New(constants.BodyWidth(), constants.BodyHeight())

	if content != "" {
		vp.SetContent(constants.ViewportContent(content))
	} else {
		vp.SetContent(constants.ViewportContent("Operation completed successfully"))
	}

	return SuccessModle{
		ReturnMc: returnMc,
		Viewport: vp,
	}
}

func GoSuccess(content string, returnFun errReturnModel) (tea.Model, tea.Cmd) {
	homeModel := makeSuccessModle(content, returnFun)
	return homeModel, nil
}

func (m SuccessModle) Init() tea.Cmd {
	return nil
}

func (m SuccessModle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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

func (m SuccessModle) View() string {
	return constants.Layout("SUCCESS", "q: Quit", m.Viewport.View())
}
