package parser

import (
	"strings"
	"testing"
)

func TestParse_ValidJSON(t *testing.T) {
	input := `{
		"model": {"display_name": "Claude Opus 4"},
		"workspace": {"current_dir": "/home/user/project"},
		"cost": {
			"total_cost_usd": 0.123,
			"total_lines_added": 50,
			"total_lines_removed": 10,
			"total_duration_ms": 60000
		},
		"context_window": {
			"context_window_size": 200000,
			"total_input_tokens": 5000,
			"total_output_tokens": 3000,
			"current_usage": {
				"input_tokens": 1000,
				"cache_creation_input_tokens": 200,
				"cache_read_input_tokens": 300
			}
		}
	}`

	data, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Model.DisplayName != "Claude Opus 4" {
		t.Errorf("model display_name = %q, want %q", data.Model.DisplayName, "Claude Opus 4")
	}
	if data.Workspace.CurrentDir != "/home/user/project" {
		t.Errorf("workspace current_dir = %q, want %q", data.Workspace.CurrentDir, "/home/user/project")
	}
	if data.Cost.TotalCostUSD != 0.123 {
		t.Errorf("cost total_cost_usd = %f, want %f", data.Cost.TotalCostUSD, 0.123)
	}
	if data.Cost.TotalLinesAdded != 50 {
		t.Errorf("cost total_lines_added = %d, want %d", data.Cost.TotalLinesAdded, 50)
	}
	if data.Cost.TotalLinesRemoved != 10 {
		t.Errorf("cost total_lines_removed = %d, want %d", data.Cost.TotalLinesRemoved, 10)
	}
	if data.Cost.TotalDurationMs != 60000 {
		t.Errorf("cost total_duration_ms = %d, want %d", data.Cost.TotalDurationMs, 60000)
	}
	if data.ContextWindow.ContextWindowSize != 200000 {
		t.Errorf("context_window_size = %d, want %d", data.ContextWindow.ContextWindowSize, 200000)
	}
	if data.ContextWindow.CurrentUsage == nil {
		t.Fatal("current_usage should not be nil")
	}
	if data.ContextWindow.CurrentUsage.InputTokens != 1000 {
		t.Errorf("input_tokens = %d, want %d", data.ContextWindow.CurrentUsage.InputTokens, 1000)
	}
}

func TestParse_EmptyInput(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestParse_NilReader(t *testing.T) {
	_, err := Parse(nil)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
}

func TestParse_MalformedJSON(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"truncated", `{"model": {"display_name": "test"`},
		{"invalid_syntax", `{model: bad}`},
		{"plain_text", `hello world`},
		{"array_instead_of_object", `[1, 2, 3]`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(strings.NewReader(tc.input))
			if err == nil {
				t.Errorf("expected error for malformed JSON %q, got nil", tc.name)
			}
		})
	}
}

func TestParse_MissingFields(t *testing.T) {
	// 只有 model 字段, 其余缺失 - JSON 解码应成功, 缺失字段为零值
	input := `{"model": {"display_name": "Sonnet"}}`
	data, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Model.DisplayName != "Sonnet" {
		t.Errorf("display_name = %q, want %q", data.Model.DisplayName, "Sonnet")
	}
	if data.Workspace.CurrentDir != "" {
		t.Errorf("current_dir should be empty, got %q", data.Workspace.CurrentDir)
	}
	if data.Cost.TotalCostUSD != 0 {
		t.Errorf("total_cost_usd should be 0, got %f", data.Cost.TotalCostUSD)
	}
	if data.ContextWindow.CurrentUsage != nil {
		t.Error("current_usage should be nil when missing")
	}
}

func TestParse_ExtraFields(t *testing.T) {
	// 包含未知字段, 应被忽略
	input := `{
		"model": {"display_name": "Haiku", "unknown_field": 42},
		"extra_top_level": true
	}`
	data, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Model.DisplayName != "Haiku" {
		t.Errorf("display_name = %q, want %q", data.Model.DisplayName, "Haiku")
	}
}

func TestParse_EmptyObject(t *testing.T) {
	input := `{}`
	data, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Model.DisplayName != "" {
		t.Errorf("display_name should be empty, got %q", data.Model.DisplayName)
	}
}
