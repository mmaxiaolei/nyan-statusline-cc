package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault_AllEnabled(t *testing.T) {
	cfg := Default()
	if !cfg.Line2Enabled {
		t.Error("Default config should have Line2Enabled=true")
	}
	for _, f := range Line1Fields {
		if !cfg.IsLine1Enabled(f.Key) {
			t.Errorf("Default config should have line1.%s enabled", f.Key)
		}
	}
	for _, f := range Line2Fields {
		if !cfg.IsLine2Enabled(f.Key) {
			t.Errorf("Default config should have line2.%s enabled", f.Key)
		}
	}
}

func TestLoad_NoFile_ReturnsDefault(t *testing.T) {
	cfg := Load(t.TempDir())
	if !cfg.Line2Enabled {
		t.Error("Load with no file should return default (Line2Enabled=true)")
	}
}

func TestLoad_InvalidJSON_ReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, configFileName), []byte("invalid"), 0644)
	cfg := Load(dir)
	if !cfg.Line2Enabled {
		t.Error("Load with invalid JSON should return default")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := Default()
	cfg.Line2Enabled = false
	cfg.Line1["model"] = false

	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded := Load(dir)
	if loaded.Line2Enabled {
		t.Error("Loaded config should have Line2Enabled=false")
	}
	if loaded.IsLine1Enabled("model") {
		t.Error("Loaded config should have line1.model disabled")
	}
	if !loaded.IsLine1Enabled("git") {
		t.Error("Loaded config should have line1.git enabled")
	}
}

func TestIsLine2Enabled_Line2Disabled(t *testing.T) {
	cfg := Default()
	cfg.Line2Enabled = false
	for _, f := range Line2Fields {
		if cfg.IsLine2Enabled(f.Key) {
			t.Errorf("IsLine2Enabled(%s) should be false when Line2Enabled=false", f.Key)
		}
	}
}

func TestIsLine1Enabled_UnknownKey(t *testing.T) {
	cfg := Default()
	if !cfg.IsLine1Enabled("nonexistent") {
		t.Error("Unknown key should default to enabled")
	}
}
