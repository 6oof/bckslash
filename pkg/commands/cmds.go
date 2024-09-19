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
			return ExecFinishedMsg{Err: errors.New("could not find editor command"), Content: ""}
		}
	}

	c := exec.Command(editorCommand, filepath) //nolint:gosec
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
		content, err := os.ReadFile("mds/HELP.md")
		if err != nil {
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Err: nil, Content: string(content)}
	}
}

func OpenProjectBcksCompose(uuid string) tea.Cmd {
	return func() tea.Msg {
		// Read the content of the HELP.md file
		content, err := os.ReadFile(path.Join("projects", uuid, "bckslash-compose.yaml"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("Bckslash compose file not found, to run the project please follow the instructins in home>help")}
			}
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Err: nil, Content: string(content)}
	}
}

func OpenProjectBcksDeployScript(uuid string) tea.Cmd {
	return func() tea.Msg {
		// Read the content of the HELP.md file
		content, err := os.ReadFile(path.Join("projects", uuid, "bckslash-deploy.sh"))
		if err != nil {

			if errors.Is(err, os.ErrNotExist) {
				return ProgramErrMsg{Err: errors.New("bckslash-deploy.sh file not found!\nThis might not be a problem if the only thing you're looking to do when deploying is starting the containers.")}
			}
			return ProgramErrMsg{Err: err}
		}
		return ExecFinishedMsg{Err: nil, Content: string(content)}
	}
}

func ShowNeofetch() tea.Cmd {
	return func() tea.Msg {
		// Start a goroutine to run the command asynchronously
		c := exec.Command("neofetch", "--off") //nolint:gosec
		out, err := c.Output()

		// Return the result as a message when the command finishes
		if err != nil {
			return ExecFinishedMsg{Err: err, Content: ""}
		}
		return ExecFinishedMsg{Err: nil, Content: string(out)}
	}
}

func ShowProjectStatus(uuid string) tea.Cmd {
	return func() tea.Msg {
		projectPath := path.Join("projects", uuid)

		// Define styles for titles and content
		titleStyle := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(constants.HighlightColor)
		contentStyle := lipgloss.NewStyle().PaddingLeft(2)

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

		// Docker Compose ps output or error

		combinedOut.WriteString(titleStyle.Render("\nDocker Compose Status\n"))
		combinedOut.WriteString("\n")
		if dockerErr != nil {
			combinedOut.WriteString(contentStyle.Render("Error running docker-compose ps:\n" + string(dockerOut) + "\n" + dockerErr.Error() + "\n"))
		} else {
			combinedOut.WriteString(contentStyle.Render(string(dockerOut) + "\n"))
		}

		// Git status output or error
		combinedOut.WriteString(titleStyle.Render("Git Status\n"))
		combinedOut.WriteString("\n")
		if gitErr != nil {
			combinedOut.WriteString(contentStyle.Render("Error running git status:\n\n" + string(gitOut) + "\n" + gitErr.Error() + "\n"))
		} else {
			combinedOut.WriteString(contentStyle.Render(string(gitOut) + "\n\n"))
		}

		// Return the combined output as a message
		return ExecFinishedMsg{Err: nil, Content: combinedOut.String()}
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

func TriggerDeploy(uuid string) tea.Cmd {

	deployType, err := helpers.DeployCheck(uuid)

	depSh := func() tea.Cmd {
		pdir := path.Join("projects", uuid)

		c := exec.Command("/bin/sh", "bckslash-deploy.sh")
		c.Dir = pdir

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = &stdoutBuf
		c.Stderr = &stderrBuf

		return tea.ExecProcess(c, func(err error) tea.Msg {

			if err != nil {
				// Handle both stderr content and the error
				if stderrBuf.Len() > 0 {
					return ExecFinishedMsg{Err: fmt.Errorf("error: %v, stderr: %s", err, stderrBuf.String())}
				}
			}

			return ExecFinishedMsg{Content: stdoutBuf.String()}
		})
	}

	depPlain := func() tea.Cmd {
		pdir := path.Join("projects", uuid)

		c := exec.Command("docker-copose", "up", "-f", "bckslash-compose.yaml", "-d")
		c.Dir = pdir

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = &stdoutBuf
		c.Stderr = &stderrBuf

		return tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				// Handle both stderr content and the error
				if stderrBuf.Len() > 0 {
					return ExecFinishedMsg{Err: fmt.Errorf("error: %v, stderr: %s", err, stderrBuf.String())}
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
	case helpers.DeployPlain:
		return depPlain()
	default:
		return nil
	}

}
