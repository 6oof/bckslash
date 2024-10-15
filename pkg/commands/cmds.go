package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/6oof/bckslash/pkg/constants"
	"github.com/6oof/bckslash/pkg/helpers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/wish"
)

func FetchProject(uuid string) tea.Cmd {
	return func() tea.Msg {
		proj, err := helpers.GetProject(uuid)
		if err != nil {
			return ProgramErrMsg{
				Err: err,
			}
		}
		return ProjectFoundMsg{Project: proj}
	}
}

func OpenEditor(filepath string) tea.Cmd {
	editorCommand := helpers.GetEditorSetting()
	if editorCommand == "" {
		return func() tea.Msg {
			return ProgramErrMsg{Err: errors.New("could not find editor command")}
		}
	}

	c := exec.Command(editorCommand, filepath) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: "Editor closed successfully"}
	})
}

func OpenHtop() tea.Cmd {
	c := exec.Command("htop", "--readonly") //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{}
	})
}

func FetchProjectGitStatus(uuid string) tea.Cmd {
	return func() tea.Msg {

		status, err := helpers.FetchProjectGitStatus(uuid)

		if err != nil {
			return ProjectViewData{GitLog: err.Error()}
		}
		return ProjectViewData{GitLog: status}
	}
}

func OpenHelpMd() tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile("public/HELP.md")
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: string(content)}
	}
}

func ReadProjectActions(uuid string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path.Join(constants.ProjectsDir, uuid, "bckslash-actions.sh"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("bckslash-actions.sh file not found! To deploy the project please follow the instructins in home>help")}
			}
			return ProgramErrMsg{Err: err}
		}

		var cntnt string = "```bash\n" + string(content) + "```"

		out, err := glamour.Render(cntnt, "tokyo-night")
		if err != nil {
			return ProgramErrMsg{Err: errors.New("Error rendering markdown")}
		}

		return ExecFinishedMsg{Content: out}
	}
}

func ReadProjectDomain(uuid string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path.Join(constants.ProjectsDir, uuid, "bckslash.caddy"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("bckslash.caddy file not found! To deploy the project please follow the instructins in home>help")}
			}
			return ProgramErrMsg{Err: err}
		}

		var cntnt string = "```caddyfile\n" + string(content) + "```"

		out, err := glamour.Render(cntnt, "tokyo-night")
		if err != nil {
			return ProgramErrMsg{Err: errors.New("Error rendering markdown")}
		}

		return ExecFinishedMsg{Content: out}
	}
}

func ShowNeofetch() tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("screenfetch") //nolint:gosec
		out, err := c.Output()

		if err != nil {

			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: string(out)}
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
		projects, err := helpers.GetProjects()
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ProjectListChangedMsg{ProjectList: projects}
	}
}

func TriggerAction(uuid, action string) tea.Cmd {
	deployType, err := helpers.DeployCheck(uuid, "projects")
	if err != nil {
		return func() tea.Msg {
			return ProgramErrMsg{Err: err}
		}
	}

	if deployType == helpers.UnDeployable {
		return func() tea.Msg {
			return ProgramErrMsg{Err: errors.New("project is undeployable, make sure you satisfy the minimum deploy requirements in home>help")}
		}
	}

	scriptPath := path.Join(constants.ProjectsDir, uuid, "bckslash-actions.sh")

	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return func() tea.Msg {
			return ProgramErrMsg{Err: fmt.Errorf("could not read script: %v", err)}
		}
	}

	if strings.Contains(string(scriptContent), "sudo") {
		return func() tea.Msg {
			return ProgramErrMsg{Err: errors.New("Error: 'sudo' is not allowed in the action script.")}
		}
	}

	// If we reach here, we are ready to execute the action.
	pdir := path.Join(constants.ProjectsDir, uuid)
	cmd := exec.Command("/bin/sh", "-c", "bckslash-actions.sh", action)
	cmd.Dir = pdir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			if stderrBuf.Len() > 0 {
				return ProgramErrMsg{Err: fmt.Errorf(" %v\nOutput: %s", err, stderrBuf.String())}
			}
			return ProgramErrMsg{Err: err} // fallback to the original error
		}
		return ExecFinishedMsg{Content: stdoutBuf.String()}
	})
}

func CommandInProject(uuid string, command string) tea.Cmd {
	// // Check if the command contains "sudo"
	// if strings.Contains(command, "sudo") {
	// 	return func() tea.Msg {
	// 		return ProgramErrMsg{Err: errors.New("Error: 'sudo' is not allowed in this shell.")}
	// 	}
	// }

	c := wish.Command(constants.WishSession, "bash", "-im")
	c.SetDir(path.Join(constants.ProjectsDir, uuid))

	cmd := tea.Exec(c, func(err error) tea.Msg {
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{}
	})

	return cmd
}
