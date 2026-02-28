// Package config 管理 nyan-statusline 的显示配置
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "nyan-config.json"

// Config 状态栏显示配置
type Config struct {
	Line2Enabled bool            `json:"line2_enabled"`
	Line1        map[string]bool `json:"line1"`
	Line2        map[string]bool `json:"line2"`
}

// Line1 字段 key 列表 (有序)
var Line1Fields = []FieldDef{
	{"model", "🤖 模型名称"},
	{"dir", "📁 项目目录"},
	{"git", "🌿 Git 分支"},
	{"context", "🌈 上下文进度"},
	{"cost", "💰 成本"},
	{"changes", "+/- 代码变更"},
	{"duration", "⏱️ 会话时长"},
	{"tokens", "📥📤 Token"},
	{"nyan", "🐱 Nyan Cat"},
	{"heartbeat", "💗 心跳动画"},
}

// Line2 字段 key 列表 (有序)
var Line2Fields = []FieldDef{
	{"codingDays", "📅 使用天数"},
	{"activeDays", "🔥 活跃天数"},
	{"streak", "⚡ 连续活跃"},
	{"sessions", "💬 会话数"},
	{"messages", "🗣️ 消息数"},
	{"todayMessages", "📈 今日统计"},
	{"peakHour", "🕐 高峰时段"},
	{"achievement", "🏆 成就徽章"},
	{"randomStatus", "🎲 随机状态"},
}

// FieldDef 字段定义
type FieldDef struct {
	Key   string
	Label string
}

// Default 返回默认配置 (全部启用)
func Default() *Config {
	c := &Config{
		Line2Enabled: true,
		Line1:        make(map[string]bool),
		Line2:        make(map[string]bool),
	}
	for _, f := range Line1Fields {
		c.Line1[f.Key] = true
	}
	for _, f := range Line2Fields {
		c.Line2[f.Key] = true
	}
	return c
}

// Load 从指定目录加载配置, 文件不存在则返回默认配置
func Load(dir string) *Config {
	path := filepath.Join(dir, configFileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return Default()
	}
	c := Default()
	if err := json.Unmarshal(raw, c); err != nil {
		return Default()
	}
	return c
}

// Save 将配置保存到指定目录
func Save(dir string, c *Config) error {
	path := filepath.Join(dir, configFileName)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// IsLine1Enabled 查询 line1 某字段是否启用
func (c *Config) IsLine1Enabled(key string) bool {
	v, ok := c.Line1[key]
	return !ok || v // 未配置的字段默认启用
}

// IsLine2Enabled 查询 line2 某字段是否启用
func (c *Config) IsLine2Enabled(key string) bool {
	if !c.Line2Enabled {
		return false
	}
	v, ok := c.Line2[key]
	return !ok || v
}
