// Package settings handles reading and writing .claude/settings.local.json.
package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const relPath = ".claude/settings.local.json"

// Load reads settings.local.json from .claude/ inside dir.
// Returns an empty map (not an error) when the file does not exist.
func Load(dir string) (map[string]interface{}, error) {
	path := filepath.Join(dir, relPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]interface{}{}, nil
		}
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Save writes data as pretty JSON to .claude/settings.local.json inside dir,
// creating the directory if needed.
func Save(dir string, data map[string]interface{}) error {
	path := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

// FilePath returns the absolute path to settings.local.json for dir.
func FilePath(dir string) string {
	return filepath.Join(dir, relPath)
}
