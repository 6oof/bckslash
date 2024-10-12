package helpers

import (
	"os"
	"path"
	"testing"

	"github.com/6oof/bckslash/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func setupTestDb(t *testing.T) string {
	testDb := "test.db"
	OpenDb(testDb)

	tempDir, _ := os.MkdirTemp("", "test-deploy-*")
	constants.ProjectsDir = tempDir

	constants.Testing = true

	t.Cleanup(func() {
		os.Remove(testDb)
		CloseDb()

		os.RemoveAll(tempDir)
	})

	return tempDir
}

func addProjectFormComandSuccess(name string) error {
	err := AddProjectFromCommand(name, "test", "https://github.com/6oof/bckslash.git", "main", "test", "test.com")

	return err
}

func TestAddProjectFormCommand(t *testing.T) {
	setupTestDb(t)

	err := addProjectFormComandSuccess("test")

	assert.NoError(t, err, "project couldn't be added")
}

func TestAddAndRemoveWithGit(t *testing.T) {
	tempDir := setupTestDb(t)
	constants.Testing = false

	err := addProjectFormComandSuccess("test")
	assert.NoError(t, err, "project couldn't be added")

	p, _ := GetProjects()
	pId := p[0].UUID

	assert.DirExists(t, path.Join(tempDir, pId), "repository folder not created")

	err = RemoveProject(pId)
	assert.NoError(t, err, "project couldn't be removed")
}

func TestAddProjectFormCommandWithGitFail(t *testing.T) {
	setupTestDb(t)
	constants.Testing = false

	err := AddProjectFromCommand("fails", "test", "https://github.com/6oof/bcfail.git", "main", "test", "test.com")

	assert.Error(t, err, "bad repository project creation didn't fail")

	p, _ := GetProjects()
	assert.Len(t, p, 0, "project was saved to the database")
}

func TestGetProjects(t *testing.T) {
	setupTestDb(t)
	addProjectFormComandSuccess("test")
	addProjectFormComandSuccess("test2")

	p, err := GetProjects()
	assert.NoError(t, err, "couldn't get projects")

	titles := []string{p[0].Title, p[1].Title}
	assert.Contains(t, titles, "test", "didn't get correct projects")
	assert.Contains(t, titles, "test2", "didn't get correct projects")
}

func TestGetProject(t *testing.T) {
	setupTestDb(t)
	addProjectFormComandSuccess("test")

	p, _ := GetProjects()
	pId := p[0].UUID

	sP, err := GetProject(pId)
	assert.NoError(t, err, "got error when getting project")

	assert.Equal(t, sP.UUID, pId, "didn't get correct projects")
}
