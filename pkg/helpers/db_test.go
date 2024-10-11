package helpers

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseOpen(t *testing.T) {
	testDb := "test.db"
	defer os.Remove(testDb)

	err := OpenDb(testDb)
	assert.NoError(t, err, "Opening a database failed")
	assert.NotNil(t, database, "Database is nil")

	err = CloseDb()
	assert.NoError(t, err, "Closing the database failed")
	assert.Equal(t, "", database.Path(), "Database is not nil")
}

func TestEnsureSettings(t *testing.T) {
	testDb := "test.db"
	defer os.Remove(testDb)
	defer CloseDb()
	OpenDb(testDb)

	err := ensureSettings()
	assert.NoError(t, err, "ensureSettings returned an error")

	err = database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		assert.NotNil(t, b, "Settings bucket not found")
		editor := b.Get([]byte("editor_command"))
		assert.Equal(t, []byte("nano"), editor, "Default editor should be 'nano'")
		return nil
	})

	assert.NoError(t, err, "View returned an error")
}

func TestEnsureProjects(t *testing.T) {
	testDb := "test.db"
	defer os.Remove(testDb)
	defer CloseDb()
	OpenDb(testDb)

	err := ensureProjects()
	assert.NoError(t, err, "ensureProjects returned an error")

	err = database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		assert.NotNil(t, b, "Projects bucket not found")
		return nil
	})

	assert.NoError(t, err, "View returned an error")
}
