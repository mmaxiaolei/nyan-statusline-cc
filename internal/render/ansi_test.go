package render

import (
	"strings"
	"testing"
)

// TestColorize_WrapsText 验证 Colorize 正确包装文本
func TestColorize_WrapsText(t *testing.T) {
	result := Colorize("hello", Red)
	if !strings.Contains(result, "hello") {
		t.Error("Colorize should contain the original text")
	}
	if !strings.HasPrefix(result, Red) {
		t.Errorf("Colorize should start with color code, got: %q", result)
	}
}

// TestColorize_ContainsReset 验证 Colorize 输出以 Reset 结尾
func TestColorize_ContainsReset(t *testing.T) {
	result := Colorize("test", Green)
	if !strings.HasSuffix(result, Reset) {
		t.Errorf("Colorize should end with Reset, got: %q", result)
	}
}

// TestColorize_Format 验证 Colorize 输出格式为 color+text+reset
func TestColorize_Format(t *testing.T) {
	result := Colorize("abc", Blue)
	expected := Blue + "abc" + Reset
	if result != expected {
		t.Errorf("Colorize format mismatch: got %q, want %q", result, expected)
	}
}

// TestColorize_EmptyText 验证空文本也能正确包装
func TestColorize_EmptyText(t *testing.T) {
	result := Colorize("", Cyan)
	expected := Cyan + Reset
	if result != expected {
		t.Errorf("Colorize empty text: got %q, want %q", result, expected)
	}
}

// TestColorize_AllColors 验证所有颜色常量都能正确包装
func TestColorize_AllColors(t *testing.T) {
	colors := map[string]string{
		"Black":   Black,
		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
	}
	for name, color := range colors {
		t.Run(name, func(t *testing.T) {
			result := Colorize("x", color)
			if !strings.HasPrefix(result, color) {
				t.Errorf("Colorize with %s should start with its color code", name)
			}
			if !strings.HasSuffix(result, Reset) {
				t.Errorf("Colorize with %s should end with Reset", name)
			}
		})
	}
}

// TestContextColor_LowUsage 验证低使用率返回绿色
func TestContextColor_LowUsage(t *testing.T) {
	if c := ContextColor(10); c != Green {
		t.Errorf("ContextColor(10) should be Green, got %q", c)
	}
	if c := ContextColor(0); c != Green {
		t.Errorf("ContextColor(0) should be Green, got %q", c)
	}
	if c := ContextColor(29.9); c != Green {
		t.Errorf("ContextColor(29.9) should be Green, got %q", c)
	}
}

// TestContextColor_MediumUsage 验证中等使用率返回黄色
func TestContextColor_MediumUsage(t *testing.T) {
	if c := ContextColor(30); c != Yellow {
		t.Errorf("ContextColor(30) should be Yellow, got %q", c)
	}
	if c := ContextColor(50); c != Yellow {
		t.Errorf("ContextColor(50) should be Yellow, got %q", c)
	}
	if c := ContextColor(79.9); c != Yellow {
		t.Errorf("ContextColor(79.9) should be Yellow, got %q", c)
	}
}

// TestContextColor_HighUsage 验证高使用率返回红色
func TestContextColor_HighUsage(t *testing.T) {
	if c := ContextColor(80); c != Red {
		t.Errorf("ContextColor(80) should be Red, got %q", c)
	}
	if c := ContextColor(100); c != Red {
		t.Errorf("ContextColor(100) should be Red, got %q", c)
	}
}
