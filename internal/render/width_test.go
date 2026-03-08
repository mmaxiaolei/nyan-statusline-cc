package render

import (
	"strings"
	"testing"
)

// TestVisualWidth_ASCII 验证纯 ASCII 字符串宽度等于字节数
func TestVisualWidth_ASCII(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"", 0},
		{"+45 -12", 7},
		{"$0.123", 6},
	}
	for _, tt := range tests {
		if got := VisualWidth(tt.input); got != tt.want {
			t.Errorf("VisualWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// TestVisualWidth_StripANSI 验证 ANSI 转义序列不计入宽度
func TestVisualWidth_StripANSI(t *testing.T) {
	// "\x1b[92mhello\x1b[0m" 视觉宽度应等于 "hello" = 5
	colored := Colorize("hello", Green)
	if got := VisualWidth(colored); got != 5 {
		t.Errorf("VisualWidth(colored) = %d, want 5", got)
	}
	// 256 色序列
	seq256 := "\x1b[38;5;220mABC\x1b[0m"
	if got := VisualWidth(seq256); got != 3 {
		t.Errorf("VisualWidth(256color) = %d, want 3", got)
	}
}

// TestVisualWidth_Emoji 验证 Emoji 宽度为 2
func TestVisualWidth_Emoji(t *testing.T) {
	// 🐱 = U+1F431, 宽度 2
	if got := VisualWidth("🐱"); got != 2 {
		t.Errorf("VisualWidth(🐱) = %d, want 2", got)
	}
	// 📁 = U+1F4C1, 宽度 2
	if got := VisualWidth("📁 abc"); got != 2+1+3 {
		t.Errorf("VisualWidth(📁 abc) = %d, want 6", got)
	}
}

// TestVisualWidth_CJK 验证 CJK 字符宽度为 2
func TestVisualWidth_CJK(t *testing.T) {
	// "天" = U+5929, 宽度 2
	if got := VisualWidth("天"); got != 2 {
		t.Errorf("VisualWidth(天) = %d, want 2", got)
	}
	// "3天" = ASCII*1 + CJK*2 = 1+2 = 3
	if got := VisualWidth("3天"); got != 3 {
		t.Errorf("VisualWidth(3天) = %d, want 3", got)
	}
}

// TestVisualWidth_Separator 验证分隔符 │ (U+2502) 宽度为 1
func TestVisualWidth_Separator(t *testing.T) {
	// │ 在 Box Drawing 区 (0x2500-0x257F), 应为 1 列
	if got := VisualWidth("│"); got != 1 {
		t.Errorf("VisualWidth(│) = %d, want 1", got)
	}
	// " │ " = 3 列
	if got := VisualWidth(" │ "); got != 3 {
		t.Errorf("VisualWidth( │ ) = %d, want 3", got)
	}
}

// TestWrapParts_NoWrapWhenFitsInWidth 验证内容未超出宽度时保持单行
func TestWrapParts_NoWrapWhenFitsInWidth(t *testing.T) {
	parts := []string{"abc", "def", "ghi"}
	sep := " │ "
	// 总宽度: 3 + 3 + 3 + 3 + 3 + 3 = 3+3+3 + 3+3 = 15 (3 parts + 2 seps = 3*3+2*3=15)
	result := wrapParts(parts, sep, 200)
	if strings.Count(result, "\n") != 0 {
		t.Errorf("wrapParts should not wrap when content fits: %q", result)
	}
	expected := "abc │ def │ ghi"
	if result != expected {
		t.Errorf("wrapParts = %q, want %q", result, expected)
	}
}

// TestWrapParts_WrapWhenExceedsWidth 验证超出宽度时自动折行
func TestWrapParts_WrapWhenExceedsWidth(t *testing.T) {
	// 每个 part 宽 5, sep 宽 3
	// 宽度 10: 第一行放 "hello"(5), 加 sep+world = 3+5=8 => 5+8=13 > 10, 折行
	parts := []string{"hello", "world", "foo"}
	sep := " │ "
	result := wrapParts(parts, sep, 10)
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Errorf("wrapParts should wrap into multiple lines for width=10, got %q", result)
	}
}

// TestWrapParts_SinglePartAlwaysFits 验证单个 part 无论多宽都独占一行不截断
func TestWrapParts_SinglePartAlwaysFits(t *testing.T) {
	parts := []string{"a very long string that exceeds any reasonable terminal width limit"}
	result := wrapParts(parts, " │ ", 5)
	if strings.Count(result, "\n") != 0 {
		t.Errorf("single part should not be wrapped: %q", result)
	}
}

// TestWrapParts_EmptyParts 验证空列表返回空字符串
func TestWrapParts_EmptyParts(t *testing.T) {
	if got := wrapParts(nil, " │ ", 80); got != "" {
		t.Errorf("wrapParts(nil) = %q, want empty", got)
	}
	if got := wrapParts([]string{}, " │ ", 80); got != "" {
		t.Errorf("wrapParts([]) = %q, want empty", got)
	}
}

// TestWrapParts_ZeroWidth 验证宽度异常时退化为单行
func TestWrapParts_ZeroWidth(t *testing.T) {
	parts := []string{"a", "b", "c"}
	result := wrapParts(parts, " │ ", 0)
	if strings.Count(result, "\n") != 0 {
		t.Errorf("wrapParts with termWidth=0 should return single line: %q", result)
	}
}

// TestWrapParts_ExactWidth 验证恰好等于宽度时不折行
func TestWrapParts_ExactWidth(t *testing.T) {
	// "ab"(2) + " │ "(3) + "cd"(2) = 7
	parts := []string{"ab", "cd"}
	result := wrapParts(parts, " │ ", 7)
	if strings.Count(result, "\n") != 0 {
		t.Errorf("wrapParts at exact width should not wrap: %q", result)
	}
}

// TestRuneWidth_CommonRanges 验证各 Unicode 区间宽度分类正确
func TestRuneWidth_CommonRanges(t *testing.T) {
	tests := []struct {
		r    rune
		want int
		name string
	}{
		{'A', 1, "ASCII letter"},
		{' ', 1, "ASCII space"},
		{'\n', 0, "newline control"},
		{'\x1b', 0, "ESC control"},
		{'│', 1, "Box drawing U+2502"},
		{'✨', 1, "Sparkles U+2728 (narrow range)"},
		{'天', 2, "CJK U+5929"},
		{'🐱', 2, "Cat emoji U+1F431"},
		{'日', 2, "CJK U+65E5"},
		{'한', 2, "Hangul U+D55C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runeWidth(tt.r); got != tt.want {
				t.Errorf("runeWidth(%q U+%04X) = %d, want %d", tt.r, tt.r, got, tt.want)
			}
		})
	}
}
