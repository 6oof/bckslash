package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/lithammer/shortuuid/v4"
)

// Mutex to ensure thread-safe access to projects
var projectMutex sync.Mutex

type Project struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Type       string `json:"type"`
}

// GetProjects reads the JSON file and unmarshals it into a slice of Project structs
func GetProjects() ([]Project, error) {

	// Open the projects file
	file, err := os.Open("bckslash_projects.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty slice if the file doesn't exist yet
			return []Project{}, nil
		}
		return nil, fmt.Errorf("unable to open projects file: %w", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read projects file: %w", err)
	}

	// Unmarshal the JSON into a slice of Project structs
	var projects []Project
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, fmt.Errorf("unable to parse projects: %w", err)
	}

	return projects, nil
}

// SaveProjects marshals the slice of Project structs and writes it to the JSON file
func SaveProjects(projects []Project) error {
	// Lock the mutex to ensure thread-safe access
	projectMutex.Lock()
	defer projectMutex.Unlock()

	// Marshal the Project slice into JSON with indentation
	bytes, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal projects: %w", err)
	}

	// Write the JSON content back to the file
	err = os.WriteFile("bckslash_projects.json", bytes, 0644)
	if err != nil {
		return fmt.Errorf("unable to write projects file: %w", err)
	}

	return nil
}

// AddProject creates a new project and appends it to the projects list
func AddProject(pro Project) error {
	// Get the current list of projects
	ap, err := GetProjects()
	if err != nil {
		return err
	}

	ap = append(ap, pro)

	// Save the updated projects list
	return SaveProjects(ap)
}

// AddProjectFromCommand clones the repository and adds the project
func AddProjectFromCommand(title, projectType, repo, branch string) error {
	pro := Project{
		Title:      title,
		UUID:       shortuuid.New(),
		Repository: repo,
		Branch:     branch,
		Type:       projectType,
	}
	// Define the path for the new project folder
	projectDir := filepath.Join("projects", pro.UUID)

	// Run the Git command to clone the repository
	c := exec.Command("git", "clone", "--depth", "1", "-b", pro.Branch, pro.Repository, projectDir)

	var stdoutBuf, stderrBuf bytes.Buffer
	c.Stdout = &stdoutBuf
	c.Stderr = &stderrBuf

	// Execute the command and capture the output
	if err := c.Run(); err != nil {
		return fmt.Errorf("git clone failed: %v\nstdout: %s\nstderr: %s\n If you're sure the repository exists, please add the Deploy key (ssh)", err, stdoutBuf.String(), stderrBuf.String())
	}

	// If cloning succeeded, proceed with adding the project
	return AddProject(pro)
}

// RemoveProject removes a project from the list by its UUID
func RemoveProject(projectUUID string) error {

	// Get the current list of projects
	projects, err := GetProjects()
	if err != nil {
		return err
	}

	// Filter out the project with the given UUID
	var newProjects []Project
	found := false
	for _, project := range projects {
		if project.UUID != projectUUID {
			newProjects = append(newProjects, project)
		} else {
			found = true
		}
	}

	// Check if the project was found
	if !found {
		return fmt.Errorf("project with UUID %s not found", projectUUID)
	}

	// Define the path for the project folder
	projectDir := filepath.Join("projects", projectUUID)
	dockerComposeFile := filepath.Join(projectDir, "docker-compose-bckslash.yml")

	// Step 0: Check if Docker Compose is running
	// Check if the docker-compose.yml file exists
	if _, err := os.Stat(dockerComposeFile); !os.IsNotExist(err) {
		psCmd := exec.Command("docker-compose", "-f", dockerComposeFile, "ps", "-q")
		psCmd.Dir = projectDir // Set the working directory to the project folder
		psOutput, err := psCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check Docker Compose status in %s: %v", projectDir, err)
		}

		if len(psOutput) > 0 {
			// Containers are running, stop them
			dockerCmd := exec.Command("docker-compose", "-f", dockerComposeFile, "down")
			dockerCmd.Dir = projectDir
			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("failed to stop Docker Compose in %s: %v", projectDir, err)
			}
		} else {
			// No containers are running
			fmt.Printf("No running Docker containers found in %s\n", projectDir)
		}
	}

	// Step 2: Remove the project folder
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to remove project folder %s: %v", projectDir, err)
	}

	// Save the updated projects list
	if err := SaveProjects(newProjects); err != nil {
		return fmt.Errorf("failed to save updated project list: %v", err)
	}

	return nil
}
