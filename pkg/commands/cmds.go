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
		// Read the content of the HELP.md file
		content, err := os.ReadFile("public/HELP.md")
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: string(content)}
	}
}

func OpenProjectBcksCompose(uuid string) tea.Cmd {
	return func() tea.Msg {
		// Read the content of the HELP.md file
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
		// Read the content of the HELP.md file
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
		// Start a goroutine to run the command asynchronously
		c := exec.Command("screenfetch") //nolint:gosec
		out, err := c.Output()

		// Return the result as a message when the command finishes
		if err != nil {

			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Content: string(out)}
	}
}

func ShowProjectStatus(uuid string) tea.Cmd {
	return func() tea.Msg {
		projectPath := path.Join(constants.ProjectsDir, uuid)

		// Define styles for titles and content
		titleStyle := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(constants.HighlightColor)
		contentStyle := lipgloss.NewStyle().PaddingLeft(2)

		// Helper function to handle command execution and formatting output
		formatOutput := func(title, output string, err error) string {
			var sb strings.Builder
			sb.WriteString(titleStyle.Render(title))
			sb.WriteString("\n\n")
			if err != nil {
				sb.WriteString(contentStyle.Render(fmt.Sprintf("Error: %s\n%s\n", err.Error(), output)))
			} else {
				sb.WriteString(contentStyle.Render(output))
			}
			sb.WriteString("\n\n") // Add uniform spacing between sections
			return sb.String()
		}

		// Execute docker-compose ps command
		cDocker := exec.Command("docker-compose", "-f", "bckslash-compose.yaml", "ps") //nolint:gosec
		cDocker.Dir = projectPath
		dockerOut, dockerErr := cDocker.CombinedOutput() // Capture both stdout and stderr

		// Execute git status command
		cGit := exec.Command("git", "status")
		cGit.Dir = projectPath
		gitOut, gitErr := cGit.CombinedOutput() // Capture both stdout and stderr

		// Combine the outputs and handle any errors
		var combinedOut strings.Builder

		// Docker Compose ps output
		combinedOut.WriteString(formatOutput("Docker Compose Status", string(dockerOut), dockerErr))

		// Git status output
		combinedOut.WriteString(formatOutput("Git Status", string(gitOut), gitErr))

		// Return the combined output as a message
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
				// Handle both stderr content and the error
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
		shell = "/bin/sh" // fallback to a default shell
	}

	// Create the command to echo the message and then start the shell
	command := fmt.Sprintf("clear && echo '\nShell in project: %s\nuse Ctrl+D or type 'exit' to exit\n' && exec %s", uuid, shell)
	c := exec.Command("sh", "-c", command) // Use "sh -c" to run the combined command
	c.Dir = path.Join(constants.ProjectsDir, uuid)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ExecFinishedMsg{}
		}
		return ExecFinishedMsg{}
	})
}
