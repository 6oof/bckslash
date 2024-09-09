package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
)

type Project struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
	Type        string `json:"type"`
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
func AddProject(title, description, projectType string) error {
	// Get the current list of projects
	projects, err := GetProjects()
	if err != nil {
		return err
	}

	// Create a new project with a UUID
	newProject := Project{
		UUID:        uuid.New().String(),
		Title:       title,
		Description: description,
		Type:        projectType,
	}

	// Append the new project to the list
	projects = append(projects, newProject)

	// Save the updated projects list
	return SaveProjects(projects)
}

// RemoveProject removes a project from the list by its UUID
func RemoveProject(projectUUID string) error {
	// Get the current list of projects
	projects, err := GetProjects()
	if err != nil {
		return err
	}

	// Filter out the project with the given UUID
	newProjects := make([]Project, 0)
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

	// Save the updated projects list
	return SaveProjects(newProjects)
}
