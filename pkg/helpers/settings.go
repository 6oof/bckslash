package helpers

import (
	"github.com/boltdb/bolt"
)

func GetEditorSetting() string {
	var ec []byte
	BcksDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		ec = b.Get([]byte("editor_command"))
		return nil
	})

	return string(ec)
}

func SetEditorSetting(editor string) error {
	err := BcksDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		b.Put([]byte("editor_command"), []byte(editor))
		return nil
	})

	if err != nil {
		return err
	}

	return nil

}
