package settings_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sapihav/clef-cli/internal/settings"
)

func TestLoad_MissingFile_ReturnsEmptyMap(t *testing.T) {
	dir := t.TempDir()
	m, err := settings.Load(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}

func TestLoad_ExistingFile_ReturnsCorrectMap(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, map[string]interface{}{
		"model":       "opus",
		"effortLevel": "xhigh",
		"otherKey":    "preserved",
	})

	m, err := settings.Load(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m["model"] != "opus" {
		t.Errorf("model: got %v, want opus", m["model"])
	}
	if m["effortLevel"] != "xhigh" {
		t.Errorf("effortLevel: got %v, want xhigh", m["effortLevel"])
	}
	if m["otherKey"] != "preserved" {
		t.Errorf("otherKey: got %v, want preserved", m["otherKey"])
	}
}

func TestSave_CreatesDirsAndFile(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{"model": "sonnet"}

	if err := settings.Save(dir, data); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	path := filepath.Join(dir, ".claude", "settings.local.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("invalid JSON written: %v", err)
	}
	if m["model"] != "sonnet" {
		t.Errorf("model: got %v, want sonnet", m["model"])
	}
}

func TestSave_PreservesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, map[string]interface{}{
		"model":       "haiku",
		"customKey":   "stays",
		"effortLevel": "low",
	})

	// Load, modify one key, save back.
	m, _ := settings.Load(dir)
	m["model"] = "opus"
	delete(m, "effortLevel")

	if err := settings.Save(dir, m); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	result, err := settings.Load(dir)
	if err != nil {
		t.Fatalf("Load after Save failed: %v", err)
	}
	if result["model"] != "opus" {
		t.Errorf("model: got %v, want opus", result["model"])
	}
	if result["customKey"] != "stays" {
		t.Errorf("customKey: got %v, want stays", result["customKey"])
	}
	if _, ok := result["effortLevel"]; ok {
		t.Errorf("effortLevel should have been removed")
	}
}

func TestReset_RemovesOnlyModelAndEffort(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, map[string]interface{}{
		"model":       "opus",
		"effortLevel": "xhigh",
		"keepThis":    true,
	})

	m, _ := settings.Load(dir)
	delete(m, "model")
	delete(m, "effortLevel")
	if err := settings.Save(dir, m); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	result, _ := settings.Load(dir)
	if _, ok := result["model"]; ok {
		t.Errorf("model should be absent after reset")
	}
	if _, ok := result["effortLevel"]; ok {
		t.Errorf("effortLevel should be absent after reset")
	}
	if result["keepThis"] != true {
		t.Errorf("keepThis should be preserved")
	}
}

// writeJSON is a test helper that writes a JSON file at .claude/settings.local.json inside dir.
func writeJSON(t *testing.T, dir string, data map[string]interface{}) {
	t.Helper()
	if err := settings.Save(dir, data); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}
}
