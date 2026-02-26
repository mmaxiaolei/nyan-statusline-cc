// Package render 提供 ANSI 终端颜色工具和状态栏渲染
package render

import "fmt"

// ANSI 颜色常量
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Black   = "\033[90m"
	Red     = "\033[91m"
	Green   = "\033[92m"
	Yellow  = "\033[93m"
	Blue    = "\033[94m"
	Magenta = "\033[95m"
	Cyan    = "\033[96m"
	White   = "\033[97m"
)

// Colorize 为文本添加 ANSI 颜色
// Parameters:
//   - text: 原始文本
//   - color: ANSI 颜色代码
//
// Return:
//   - string: 带颜色的文本
func Colorize(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// ContextColor 根据上下文使用率返回对应颜色
// Parameters:
//   - percent: 使用百分比
//
// Return:
//   - string: ANSI 颜色代码
func ContextColor(percent float64) string {
	if percent < 30 {
		return Green
	} else if percent < 80 {
		return Yellow
	}
	return Red
}
