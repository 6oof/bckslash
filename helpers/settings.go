package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type BckslashSettings struct {
	EditorCommand string `json:"editor_command"`
}

func GetSettings() (*BckslashSettings, error) {
	// Open the settings file
	file, err := os.Open("bckslash_settings.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open settings file: %w", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read settings file: %w", err)
	}

	// Unmarshal the JSON into a Settings struct
	var settings BckslashSettings
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		return nil, fmt.Errorf("unable to parse settings: %w", err)
	}

	return &settings, nil
}
