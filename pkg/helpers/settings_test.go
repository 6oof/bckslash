package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditorSettins(t *testing.T) {
	testDb := "test.db"
	defer os.Remove(testDb)
	defer CloseDb()
	OpenDb(testDb)

	editor := GetEditorSetting()

	assert.Equal(t, "nano", editor, "Default editor should be 'nano'")

}

func TestSetEditorSettins(t *testing.T) {
	testDb := "test.db"
	defer os.Remove(testDb)
	defer CloseDb()
	OpenDb(testDb)

	editor := GetEditorSetting()
	assert.Equal(t, "nano", editor, "Default editor should be 'nano'")

	err := SetEditorSetting("vim")
	assert.NoError(t, err, "Unable to set editor")

	newEditor := GetEditorSetting()
	assert.Equal(t, "vim", newEditor, "Editor should be 'vim'")
}
