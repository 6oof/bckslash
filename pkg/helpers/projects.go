package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/6oof/bckslash/pkg/constants"
	"github.com/boltdb/bolt"
	"github.com/lithammer/shortuuid/v4"
)

type Project struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Type       string `json:"type"`
}

func GetProjects() ([]Project, error) {
	var projects []Project

	err := BcksDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return fmt.Errorf("projects bucket not found")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var curp Project
			err := json.Unmarshal(v, &curp)
			if err != nil {
				return err
			}

			projects = append(projects, curp)
		}
		return nil
	})

	return projects, err
}

func GetProject(uuid string) (Project, error) {
	var project Project

	err := BcksDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return errors.New("projects bucket not found")
		}

		p := b.Get([]byte(uuid))
		if p == nil {
			return errors.New("project with selected id not found")
		}

		// Pass project by reference
		err := json.Unmarshal(p, &project)
		return err
	})

	return project, err
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

	if err := BcksDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		projectId := b.Get([]byte(pro.UUID))
		if projectId != nil {
			return errors.New("UUID collision, don't worry pelase re-try")
		}
		return nil
	}); err != nil {
		return err
	}

	if !constants.Testing {
		// Define the path for the new project folder
		projectDir := filepath.Join(constants.ProjectsDir, pro.UUID)

		// Run the Git command to clone the repository
		c := exec.Command("git", "clone", "--depth", "1", "-b", pro.Branch, pro.Repository, projectDir)

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = &stdoutBuf
		c.Stderr = &stderrBuf

		// Execute the command and capture the output
		if err := c.Run(); err != nil {
			return fmt.Errorf("git clone failed: %v\nstdout: %s\nstderr: %s\n If you're sure the repository exists, please add the Deploy key (ssh)", err, stdoutBuf.String(), stderrBuf.String())
		}

		_ = resolveEnvOnCreate(pro.UUID)

	}

	if err := BcksDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))

		pj, err := json.Marshal(pro)
		if err != nil {
			return err
		}

		err = b.Put([]byte(pro.UUID), pj)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func resolveEnvOnCreate(uuid string) error {
	envPath := path.Join(constants.ProjectsDir, uuid, ".env")
	file, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			fileExample, err := os.Open(path.Join(constants.ProjectsDir, uuid, ".env.example"))
			if err != nil {
				return fmt.Errorf("unable to open .env.example: %w", err)
			}
			defer fileExample.Close()

			bytes, err := io.ReadAll(fileExample)
			if err != nil {
				return fmt.Errorf("unable to read .env.example: %w", err)
			}

			err = os.WriteFile(envPath, bytes, 0664)
			if err != nil {
				return fmt.Errorf("unable to write .env file: %w", err)
			}
			return nil
		}
		return fmt.Errorf("unable to open .env file: %w", err)
	}
	defer file.Close()

	return nil
}

// RemoveProject removes a project from the list by its UUID
func RemoveProject(uuid string) error {

	p, err := GetProject(uuid)

	if err != nil {
		return err
	}

	// Check if the project was found
	if p.UUID == "" {
		return fmt.Errorf("project with UUID %s not found", uuid)
	}

	// Define the path for the project folder
	projectDir := filepath.Join(constants.ProjectsDir, uuid)

	// Remove the project folder
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to remove project folder %s: %v", projectDir, err)
	}

	if err = BcksDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if err := b.Delete([]byte(uuid)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func FetchProjectGitStatus(uuid string) (string, error) {

	projectDir := filepath.Join(constants.ProjectsDir, uuid)

	psCmd := exec.Command("git", "--no-pager", "log", "-1", "--format=%h %cd", "--date=iso")
	psCmd.Dir = projectDir
	psOutput, err := psCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to check git log in %s: %v", projectDir, err)
	}

	return string(psOutput), nil
}
