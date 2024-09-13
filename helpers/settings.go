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

func makeSettingsStore() error {
	settings := &BckslashSettings{
		EditorCommand: "nano",
	}

	bytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal settings: %w", err)
	}

	// Write the JSON content back to the file
	err = os.WriteFile("bckslash_settings.json", bytes, 0644)
	if err != nil {
		return fmt.Errorf("unable to write settings file: %w", err)
	}

	return nil
}

// GetSettings reads the JSON file and unmarshals it into a BckslashSettings struct
func GetSettings() (*BckslashSettings, error) {
	// Open the settings file
	file, err := os.Open("bckslash_settings.json")
	if err != nil {
		if os.IsNotExist(err) {
			err := makeSettingsStore()
			// Return an empty slice if the file doesn't exist yet
			return &BckslashSettings{}, err
		}

		return nil, fmt.Errorf("unable to open settings file: %w", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read settings file: %w", err)
	}

	// Unmarshal the JSON into a BckslashSettings struct
	var settings BckslashSettings
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		return nil, fmt.Errorf("unable to parse settings: %w", err)
	}

	return &settings, nil
}

// SaveSettings marshals the BckslashSettings struct and writes it to the JSON file
func SaveSettings(settings *BckslashSettings) error {
	// Marshal the BckslashSettings struct into JSON with indentation
	bytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal settings: %w", err)
	}

	// Write the JSON content back to the file
	err = os.WriteFile("bckslash_settings.json", bytes, 0644)
	if err != nil {
		return fmt.Errorf("unable to write settings file: %w", err)
	}

	return nil
}
