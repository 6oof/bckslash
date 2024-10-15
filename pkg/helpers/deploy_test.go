package helpers

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployCheck(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "test-deploy-*")
	defer os.RemoveAll(tempDir)

	testUUID := "test-uuid-alksjd-asdklf"

	// Test the case where the 'bcks-deploy.sh' file does not exist
	deployType, err := DeployCheck(testUUID, "test-projects")
	assert.NoError(t, err, "DeployCheck returned an unexpected error")
	assert.Equal(t, UnDeployable, deployType, "Expected DeploySh when the file exists")

	projectDir := path.Join(tempDir, testUUID)
	os.MkdirAll(projectDir, 0755)

	deployFile := path.Join(projectDir, "bckslash-actions.sh")
	os.Create(deployFile)

	// Run DeployCheck again with the file now in place
	deployType, err = DeployCheck(testUUID, tempDir)
	assert.NoError(t, err, "DeployCheck returned an unexpected error after file creation")
	assert.Equal(t, DeploySh, deployType, "Expected DeploySh when the deploy script exists")
}
