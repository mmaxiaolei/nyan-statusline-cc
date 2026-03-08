package render

import (
	"strings"
	"testing"

	"github.com/nyan-statusline-cc/internal/model"
)

// newTestSessionData 构造测试用的 SessionData
func newTestSessionData() *model.SessionData {
	return &model.SessionData{
		Model: model.ModelInfo{
			DisplayName: "claude-opus-4",
		},
		Workspace: model.WorkspaceInfo{
			CurrentDir: "/home/user/project",
		},
		Cost: model.CostInfo{
			TotalCostUSD:      0.15,
			TotalLinesAdded:   42,
			TotalLinesRemoved: 7,
			TotalDurationMs:   180000,
		},
		ContextWindow: model.ContextWindow{
			ContextWindowSize: 200000,
			TotalInputTokens:  50000,
			TotalOutputTokens: 10000,
			CurrentUsage: &model.UsageDetail{
				InputTokens:              30000,
				CacheCreationInputTokens: 5000,
				CacheReadInputTokens:     10000,
			},
		},
	}
}

// TestRender_NotEmpty 验证 Render 输出非空
func TestRender_NotEmpty(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if result == "" {
		t.Error("Render() should not return empty string")
	}
}

// TestRender_ContainsSeparator 验证输出包含分隔符
func TestRender_ContainsSeparator(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "│") {
		t.Error("Render output should contain separator '│'")
	}
}

// TestRender_ContainsModelEmoji 验证输出包含模型 emoji
func TestRender_ContainsModelEmoji(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "👾") {
		t.Error("Render output should contain model emoji '👾'")
	}
}

// TestRender_ContainsModelName 验证输出包含模型名称
func TestRender_ContainsModelName(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "claude-opus-4") {
		t.Error("Render output should contain model name")
	}
}

// TestRender_ContainsDirEmoji 验证输出包含目录 emoji
func TestRender_ContainsDirEmoji(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "🗂️") {
		t.Error("Render output should contain directory emoji '🗂️'")
	}
}

// TestRender_ContainsDirName 验证输出包含目录名 (basename)
func TestRender_ContainsDirName(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "project") {
		t.Error("Render output should contain directory basename 'project'")
	}
}

// TestRender_ContainsCost 验证输出包含成本信息
func TestRender_ContainsCost(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "💰") {
		t.Error("Render output should contain cost emoji '💰'")
	}
	if !strings.Contains(result, "$") {
		t.Error("Render output should contain dollar sign")
	}
}

// TestRender_ContainsCodeChanges 验证输出包含代码变更
func TestRender_ContainsCodeChanges(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "+42") {
		t.Error("Render output should contain added lines '+42'")
	}
	if !strings.Contains(result, "-7") {
		t.Error("Render output should contain removed lines '-7'")
	}
}

// TestRender_ContainsProgressBar 验证输出包含进度条 (彩虹色填充块)
func TestRender_ContainsProgressBar(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	// 进度条使用 "█" 填充
	if !strings.Contains(result, "█") {
		t.Error("Render output should contain progress bar filled blocks '█'")
	}
}

// TestRender_ContainsContextPercent 验证输出包含上下文使用百分比
func TestRender_ContainsContextPercent(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "%") {
		t.Error("Render output should contain context usage percentage")
	}
}

// TestRender_ContainsTokenStats 验证输出包含 Token 统计
func TestRender_ContainsTokenStats(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "📥") {
		t.Error("Render output should contain input token emoji '📥'")
	}
	if !strings.Contains(result, "📤") {
		t.Error("Render output should contain output token emoji '📤'")
	}
}

// TestRender_ContainsDuration 验证输出包含会话时长
func TestRender_ContainsDuration(t *testing.T) {
	data := newTestSessionData()
	result := Render(data)
	if !strings.Contains(result, "⏱️") {
		t.Error("Render output should contain duration emoji '⏱️'")
	}
}

// TestRender_DefaultModelName 验证模型名为空时使用 "Unknown"
func TestRender_DefaultModelName(t *testing.T) {
	data := newTestSessionData()
	data.Model.DisplayName = ""
	result := Render(data)
	if !strings.Contains(result, "Unknown") {
		t.Error("Render output should contain 'Unknown' when model name is empty")
	}
}

// TestCalcContextPercent 验证上下文使用百分比计算
func TestCalcContextPercent(t *testing.T) {
	data := newTestSessionData()
	// (30000 + 5000 + 10000) / 200000 * 100 = 22.5%
	got := calcContextPercent(data)
	if got != 22.5 {
		t.Errorf("calcContextPercent() = %f, want 22.5", got)
	}
}

// TestCalcContextPercent_ZeroWindow 验证窗口大小为 0 时返回 0
func TestCalcContextPercent_ZeroWindow(t *testing.T) {
	data := newTestSessionData()
	data.ContextWindow.ContextWindowSize = 0
	got := calcContextPercent(data)
	if got != 0 {
		t.Errorf("calcContextPercent() with zero window = %f, want 0", got)
	}
}

// TestCalcContextPercent_NilUsage 验证 CurrentUsage 为 nil 时返回 0
func TestCalcContextPercent_NilUsage(t *testing.T) {
	data := newTestSessionData()
	data.ContextWindow.CurrentUsage = nil
	got := calcContextPercent(data)
	if got != 0 {
		t.Errorf("calcContextPercent() with nil usage = %f, want 0", got)
	}
}

// TestPeakHourEmoji 验证不同时段返回正确的 emoji
func TestPeakHourEmoji(t *testing.T) {
	tests := []struct {
		hour int
		want string
	}{
		{0, "🌙"},
		{3, "🌙"},
		{4, "🌙"},
		{5, "🌅"},
		{8, "🌅"},
		{11, "🌅"},
		{12, "☀️"},
		{15, "☀️"},
		{17, "☀️"},
		{18, "🌆"},
		{20, "🌆"},
		{21, "🌆"},
		{22, "🌙"},
		{23, "🌙"},
	}
	for _, tt := range tests {
		got := peakHourEmoji(tt.hour)
		if got != tt.want {
			t.Errorf("peakHourEmoji(%d) = %q, want %q", tt.hour, got, tt.want)
		}
	}
}

// TestGetAchievement 验证成就徽章逻辑
func TestGetAchievement(t *testing.T) {
	tests := []struct {
		name string
		info *model.StatsInfo
		want string
	}{
		{"1000 messages", &model.StatsInfo{TotalMessages: 1000}, "🏆 千言万语"},
		{"500 messages", &model.StatsInfo{TotalMessages: 500}, "🥇 消息达人"},
		{"100 messages", &model.StatsInfo{TotalMessages: 100}, "🥈 话唠新星"},
		{"100 sessions", &model.StatsInfo{TotalSessions: 100}, "👑 会话之王"},
		{"50 sessions", &model.StatsInfo{TotalSessions: 50}, "⭐ 会话专家"},
		{"30 streak", &model.StatsInfo{Streak: 30}, "🔥 月度坚持"},
		{"7 streak", &model.StatsInfo{Streak: 7}, "💪 周度坚持"},
		{"3 streak", &model.StatsInfo{Streak: 3}, "✊ 三连击"},
		{"30 active days", &model.StatsInfo{ActiveDays: 30}, "🎖️ 老用户"},
		{"no achievement", &model.StatsInfo{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAchievement(tt.info)
			if got != tt.want {
				t.Errorf("getAchievement() = %q, want %q", got, tt.want)
			}
		})
	}
}
