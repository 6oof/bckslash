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

type ProgramErrMsg struct {
	Err error
}

type ProjectListChangedMsg struct {
	ProjectList []helpers.Project
}

type ReturnHomeMsg struct{}

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

func OpenGlowHelp() tea.Cmd {
	c := exec.Command("glow", "mds") //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ExecFinishedMsg{Err: err, Content: ""}
	})
}

func ShowNeofetch() tea.Cmd {
	return func() tea.Msg {
		// Start a goroutine to run the command asynchronously
		c := exec.Command("neofetch") //nolint:gosec
		out, err := c.Output()

		// Return the result as a message when the command finishes
		if err != nil {
			return ExecFinishedMsg{Err: err, Content: ""}
		}
		return ExecFinishedMsg{Err: nil, Content: string(out)}
	}
}

func AddProjectCommand(title, projectType, repo, branch string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.AddProjectFromCommand(title, projectType, repo, branch)
		if err != nil {
			return ProgramErrMsg{Err: err}
		}

		projects, err := helpers.GetProjects()
		if err != nil {
			return ProgramErrMsg{Err: err}
		}

		return ProjectListChangedMsg{
			ProjectList: projects,
		}
	}
}

func DeleteProjectCommand(uuid string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.RemoveProject(uuid)
		if err != nil {
			return ProgramErrMsg{Err: err}
		}

		projects, err := helpers.GetProjects()
		if err != nil {
			return ProgramErrMsg{Err: err}
		}

		return ProjectListChangedMsg{
			ProjectList: projects,
		}
	}

}

func LoadProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		// Load projects (e.g., from a file or database)
		projects, err := helpers.GetProjects()
		if err != nil {
			// Handle error if needed
			return ProgramErrMsg{Err: err}
		}
		// Return a message that the project list has been loaded
		return ProjectListChangedMsg{ProjectList: projects}
	}
}
