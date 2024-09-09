package commands

import (
	"lg/helpers"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type ExecFinishedMsg struct {
	Err     error
	Content string
}

type ExecStartMsg struct{}

func OpenEditor(filepath string) tea.Cmd {
	settings, err := helpers.GetSettings()
	if err != nil {
		return func() tea.Msg {
			return ExecFinishedMsg{Err: err, Content: ""}
		}
	}

	c := exec.Command(settings.EditorCommand, filepath) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ExecFinishedMsg{Err: err, Content: ""}
	})
}

func OpenBTM() tea.Cmd {
	c := exec.Command("btm", "--theme", "nord") //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ExecFinishedMsg{Err: err, Content: ""}
	})
}

func ShowNeofetch() tea.Cmd {
	c := exec.Command("nerdfetch", "-e") //nolint:gosec
	out, err := c.Output()
	if err != nil {
		return func() tea.Msg {
			return ExecFinishedMsg{Err: err, Content: ""}
		}
	} else {
		return func() tea.Msg {
			return ExecFinishedMsg{Err: nil, Content: string(out)}
		}
	}

}
