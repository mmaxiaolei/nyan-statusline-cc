// Package parser 负责解析 Claude Code 通过 stdin 传入的 JSON 数据
package parser

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nyan-statusline-cc/internal/model"
)

// Parse 从 reader 中读取并解析 JSON 数据为 SessionData
// Parameters:
//   - r: 输入源 (通常为 os.Stdin)
//
// Return:
//   - *model.SessionData: 解析后的会话数据
//   - error: 解析错误
func Parse(r io.Reader) (*model.SessionData, error) {
	if r == nil {
		return nil, errors.New("reader must not be nil")
	}
	var data model.SessionData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
