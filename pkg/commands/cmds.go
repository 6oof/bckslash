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
	"github.com/charmbracelet/lipgloss"
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

func OpenBTM() tea.Cmd {
	c := exec.Command("htop", "--readonly") //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{}
	})
}

func FetchProjectData(uuid string) tea.Cmd {
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

func OpenProjectBcksCompose(uuid string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path.Join(constants.ProjectsDir, uuid, "bckslash-compose.yaml"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("Bckslash compose file not found! To deploy the project please follow the instructins in home>help")}
			}
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: string(content)}
	}
}

func OpenProjectBcksDeployScript(uuid string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path.Join(constants.ProjectsDir, uuid, "bckslash-deploy.sh"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("bckslash-deploy.sh file not found! To deploy the project please follow the instructins in home>help")}
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

func ShowProjectStatus(uuid string) tea.Cmd {
	return func() tea.Msg {
		projectPath := path.Join(constants.ProjectsDir, uuid)

		titleStyle := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(constants.HighlightColor)
		contentStyle := lipgloss.NewStyle().PaddingLeft(2)

		formatOutput := func(title, output string, err error) string {
			var sb strings.Builder
			sb.WriteString(titleStyle.Render(title))
			sb.WriteString("\n\n")
			if err != nil {
				sb.WriteString(contentStyle.Render(fmt.Sprintf("Error: %s\n%s\n", err.Error(), output)))
			} else {
				sb.WriteString(contentStyle.Render(output))
			}
			sb.WriteString("\n\n")
			return sb.String()
		}

		cDocker := exec.Command("docker-compose", "-f", "bckslash-compose.yaml", "ps")
		cDocker.Dir = projectPath
		dockerOut, dockerErr := cDocker.CombinedOutput()
		cGit := exec.Command("git", "status")
		cGit.Dir = projectPath
		gitOut, gitErr := cGit.CombinedOutput()
		var combinedOut strings.Builder

		combinedOut.WriteString(formatOutput("Docker Compose Status", string(dockerOut), dockerErr))

		combinedOut.WriteString(formatOutput("Git Status", string(gitOut), gitErr))

		return ExecFinishedMsg{Content: combinedOut.String()}
	}
}

func AddProjectCommand(title, projectType, repo, branch, serviceName, domain string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.AddProjectFromCommand(title, projectType, repo, branch, serviceName, domain)
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

func TriggerDeploy(uuid string) tea.Cmd {

	deployType, err := helpers.DeployCheck(uuid, "projects")

	depSh := func() tea.Cmd {
		pdir := path.Join(constants.ProjectsDir, uuid)

		c := exec.Command("/bin/sh", "bckslash-deploy.sh")
		c.Dir = pdir

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = &stdoutBuf
		c.Stderr = &stderrBuf

		return tea.ExecProcess(c, func(err error) tea.Msg {

			if err != nil {
				if stderrBuf.Len() > 0 {
					return ProgramErrMsg{Err: fmt.Errorf("error: %v, stderr: %s", err, stderrBuf.String())}
				}
			}

			return ExecFinishedMsg{Content: stdoutBuf.String()}
		})
	}

	switch deployType {
	case helpers.UnDeployable:
		return func() tea.Msg { return ProgramErrMsg{Err: err} }
	case helpers.DeploySh:
		return depSh()
	default:
		return func() tea.Msg {
			return ProgramErrMsg{Err: errors.New("Project failed the minimum requirements for deployment")}
		}
	}

}

func ShellInProject(uuid string) tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	command := fmt.Sprintf("clear && echo '\nShell in project: %s\nuse Ctrl+D or type 'exit' to exit\n' && exec %s", uuid, shell)
	c := exec.Command("sh", "-c", command)
	c.Dir = path.Join(constants.ProjectsDir, uuid)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ExecFinishedMsg{}
		}
		return ExecFinishedMsg{}
	})
}
