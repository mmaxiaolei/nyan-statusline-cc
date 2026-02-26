// Package formatter 负责格式化状态栏中的各类数据显示
package formatter

import "fmt"

// FormatCost 格式化费用显示
// Parameters:
//   - cost: 费用金额 (USD)
//
// Return:
//   - string: 格式化后的费用字符串
func FormatCost(cost float64) string {
	if cost <= 0 {
		return "$0.0000"
	}
	if cost < 0.01 {
		return fmt.Sprintf("$%.4f", cost)
	} else if cost < 1 {
		return fmt.Sprintf("$%.3f", cost)
	}
	return fmt.Sprintf("$%.2f", cost)
}

// FormatDuration 格式化会话时长
// Parameters:
//   - ms: 毫秒数
//
// Return:
//   - string: 可读的时长字符串 (如 "2m30s")
func FormatDuration(ms int64) string {
	if ms <= 0 {
		return "0s"
	}
	seconds := ms / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		secs := seconds % 60
		return fmt.Sprintf("%dm%ds", minutes, secs)
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// FormatTokens 格式化 token 数量
// Parameters:
//   - tokens: token 数量
//
// Return:
//   - string: 简化后的数量字符串 (如 "50k")
func FormatTokens(tokens int64) string {
	if tokens <= 0 {
		return "0"
	}
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	} else if tokens < 10000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%dk", tokens/1000)
}
