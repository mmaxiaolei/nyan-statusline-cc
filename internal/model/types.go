// Package model 定义 Claude Code 状态栏的数据模型
package model

// SessionData 表示 Claude Code 通过 stdin 传入的完整会话数据
type SessionData struct {
	Model         ModelInfo      `json:"model"`
	Workspace     WorkspaceInfo  `json:"workspace"`
	Cost          CostInfo       `json:"cost"`
	ContextWindow ContextWindow  `json:"context_window"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	DisplayName string `json:"display_name"`
}

// WorkspaceInfo 工作区信息
type WorkspaceInfo struct {
	CurrentDir string `json:"current_dir"`
}

// CostInfo 成本和代码变更信息
type CostInfo struct {
	TotalCostUSD     float64 `json:"total_cost_usd"`
	TotalLinesAdded  int     `json:"total_lines_added"`
	TotalLinesRemoved int    `json:"total_lines_removed"`
	TotalDurationMs  int64   `json:"total_duration_ms"`
}

// ContextWindow 上下文窗口信息
type ContextWindow struct {
	ContextWindowSize int64        `json:"context_window_size"`
	TotalInputTokens  int64        `json:"total_input_tokens"`
	TotalOutputTokens int64        `json:"total_output_tokens"`
	CurrentUsage      *UsageDetail `json:"current_usage"`
}

// UsageDetail 当前上下文使用详情
type UsageDetail struct {
	InputTokens              int64 `json:"input_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}
