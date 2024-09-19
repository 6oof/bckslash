package helpers

import (
	"fmt"

	"github.com/boltdb/bolt"
)

var database *bolt.DB

func OpenDb() error {
	db, err := bolt.Open("bckslash.db", 0600, nil)

	if err != nil {
		return err
	}

	database = db

	if err := ensureSettings(); err != nil {
		return err
	}

	if err := ensureProjects(); err != nil {
		return err
	}

	return nil
}

func CloseDb() {
	database.Close()
}

func ensureSettings() error {
	database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("settings"))
		ec := b.Get([]byte("editor_command"))
		if ec == nil {
			if err := b.Put([]byte("editor_command"), []byte("nano")); err != nil {
				return err
			}
		}
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return nil
}

func ensureProjects() error {
	database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("projects"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return nil
}

func BcksDb() *bolt.DB {
	return database
}
