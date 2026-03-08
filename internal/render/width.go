package render

import (
	"regexp"
	"strings"
)

// ansiEscapeRe 匹配所有标准 CSI ANSI 转义序列 (包括颜色、256色等)
var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// sepVisualWidth 是分隔符 " │ " 的视觉列数
const sepVisualWidth = 3

// VisualWidth 计算字符串在终端中占用的视觉列数.
// 会先剥离所有 ANSI 转义序列, 再按 Unicode 宽度规则统计.
func VisualWidth(s string) int {
	clean := ansiEscapeRe.ReplaceAllString(s, "")
	w := 0
	for _, r := range clean {
		w += runeWidth(r)
	}
	return w
}

// runeWidth 返回单个 Unicode 字符在终端中的视觉列数 (1 或 2).
// 覆盖所有常见宽字符范围: CJK、Hangul、全角假名、Emoji 等.
func runeWidth(r rune) int {
	switch {
	case r < 0x20:
		// 控制字符不占列宽
		return 0
	case r < 0x1100:
		// 普通 ASCII、Latin、希腊语、西里尔等: 1 列
		return 1
	case r <= 0x115F:
		// Hangul Jamo
		return 2
	case r < 0x2E80:
		return 1
	case r <= 0x303E:
		// CJK 部首、康熙部首等
		return 2
	case r < 0x3040:
		return 1
	case r <= 0xA4CF:
		// 平假名、片假名、汉字、韩文字母等
		return 2
	case r < 0xAC00:
		return 1
	case r <= 0xD7AF:
		// 韩文音节
		return 2
	case r < 0xF900:
		return 1
	case r <= 0xFAFF:
		// CJK 兼容汉字
		return 2
	case r < 0xFE10:
		return 1
	case r <= 0xFE19:
		// 竖排形式
		return 2
	case r < 0xFE30:
		return 1
	case r <= 0xFE4F:
		// CJK 兼容形式
		return 2
	case r < 0xFF01:
		return 1
	case r <= 0xFF60:
		// 全角 ASCII 及标点
		return 2
	case r < 0xFFE0:
		return 1
	case r <= 0xFFE6:
		// 全角符号
		return 2
	case r >= 0x1B000:
		// Emoji、象形文字及其他补充平面宽字符
		return 2
	default:
		return 1
	}
}

// wrapParts 将 parts 按终端宽度自动换行, 超出 termWidth 时起新行.
// 每个 part 之间以 sep 连接, 行间以 "\n" 分隔.
// 单个 part 超出终端宽度时不强制截断, 独占一行.
func wrapParts(parts []string, sep string, termWidth int) string {
	if len(parts) == 0 {
		return ""
	}
	// 终端宽度异常时退化为单行
	if termWidth <= 0 {
		return strings.Join(parts, sep)
	}

	var lines [][]string
	currentLine := []string{}
	currentWidth := 0

	for _, part := range parts {
		w := VisualWidth(part)
		if len(currentLine) == 0 {
			// 新行的第一个 part, 直接加入
			currentLine = append(currentLine, part)
			currentWidth = w
		} else {
			needed := sepVisualWidth + w
			if currentWidth+needed > termWidth {
				// 超出宽度, 折行
				lines = append(lines, currentLine)
				currentLine = []string{part}
				currentWidth = w
			} else {
				currentLine = append(currentLine, part)
				currentWidth += needed
			}
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	rendered := make([]string, len(lines))
	for i, line := range lines {
		rendered[i] = strings.Join(line, sep)
	}
	return strings.Join(rendered, "\n")
}
