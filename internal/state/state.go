// Package state 管理 Claude Code 的处理状态 (处理中/已完成)
//
// 工作原理:
//   - hooks 调用 `nyan-statusline --state processing/completed` 写入状态
//   - statusline 渲染时读取状态文件判断当前状态
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const stateFileName = "nyan-state.json"

// Status 表示 Claude Code 的处理状态
const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

// stateData 状态文件结构
type stateData struct {
	Status string `json:"status"`
}

// SetStatus 将指定状态写入状态文件
// Parameters:
//   - binaryDir: 二进制文件所在目录 (状态文件同目录)
//   - status: 状态值 (processing/completed)
//
// Return:
//   - error: 错误信息
func SetStatus(binaryDir, status string) error {
	statePath := filepath.Join(binaryDir, stateFileName)
	data, err := json.Marshal(stateData{Status: status})
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0644)
}

// IsProcessing 读取状态文件, 判断 Claude Code 是否正在处理中
// Parameters:
//   - binaryDir: 二进制文件所在目录 (状态文件同目录)
//
// Return:
//   - bool: true 表示正在处理, false 表示已完成
func IsProcessing(binaryDir string) bool {
	statePath := filepath.Join(binaryDir, stateFileName)

	raw, err := os.ReadFile(statePath)
	if err != nil {
		// 无状态文件, 默认处理中
		return true
	}

	var s stateData
	if err := json.Unmarshal(raw, &s); err != nil {
		return true
	}

	return s.Status != StatusCompleted
}
