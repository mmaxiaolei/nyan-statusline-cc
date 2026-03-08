package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyan-statusline-cc/internal/animation"
	"github.com/nyan-statusline-cc/internal/config"
	"github.com/nyan-statusline-cc/internal/formatter"
	"github.com/nyan-statusline-cc/internal/git"
	"github.com/nyan-statusline-cc/internal/model"
	"github.com/nyan-statusline-cc/internal/state"
	"github.com/nyan-statusline-cc/internal/stats"
)

// Render 将会话数据渲染为状态栏输出字符串
// Parameters:
//   - data: Claude Code 会话数据
//
// Return:
//   - string: 完整的状态栏输出 (可能包含多行)
func Render(data *model.SessionData) string {
	sep := Colorize(" │ ", Black)

	// 加载配置
	var cfg *config.Config
	if execPath, err := os.Executable(); err == nil {
		cfg = config.Load(filepath.Dir(execPath))
	} else {
		cfg = config.Default()
	}

	line1 := renderLine1(data, sep, cfg)

	if cfg.Line2Enabled {
		if line2 := renderLine2(sep, cfg); line2 != "" {
			return line1 + "\n" + line2
		}
	}
	return line1
}

// renderLine1 渲染第一行: 模型、目录、Git、进度条、成本、变更、时长、Token、Nyan Cat、心跳
func renderLine1(data *model.SessionData, sep string, cfg *config.Config) string {
	var parts []string

	// 模型名称
	if cfg.IsLine1Enabled("model") {
		modelName := data.Model.DisplayName
		if modelName == "" {
			modelName = "Unknown"
		}
		parts = append(parts, Colorize(Bold+"👾 "+modelName, Magenta))
	}

	// 当前目录
	if cfg.IsLine1Enabled("dir") {
		dir := filepath.Base(data.Workspace.CurrentDir)
		if dir != "" {
			parts = append(parts, Colorize("🗂️ "+dir, Cyan))
		}
	}

	// Git 分支
	if cfg.IsLine1Enabled("git") {
		if gitInfo, _ := git.GetInfo(); gitInfo != nil {
			color := Green
			status := ""
			if gitInfo.HasChanges {
				color = Yellow
				status = "*"
			}
			parts = append(parts, Colorize("🌿 "+gitInfo.Branch+status, color))
		}
	}

	// 上下文使用率 + 彩虹进度条
	if cfg.IsLine1Enabled("context") {
		ctxPercent := calcContextPercent(data)
		bar := animation.RainbowProgressBar(ctxPercent, 10)
		ctxColor := ContextColor(ctxPercent)
		parts = append(parts, fmt.Sprintf("%s %s%.1f%%%s", bar, ctxColor, ctxPercent, Reset))
	}

	// 成本
	if cfg.IsLine1Enabled("cost") && data.Cost.TotalCostUSD > 0 {
		parts = append(parts, Colorize("💰 "+formatter.FormatCost(data.Cost.TotalCostUSD), Yellow))
	}

	// 代码变更
	if cfg.IsLine1Enabled("changes") && (data.Cost.TotalLinesAdded > 0 || data.Cost.TotalLinesRemoved > 0) {
		var changes []string
		if data.Cost.TotalLinesAdded > 0 {
			changes = append(changes, Colorize(fmt.Sprintf("+%d", data.Cost.TotalLinesAdded), Green))
		}
		if data.Cost.TotalLinesRemoved > 0 {
			changes = append(changes, Colorize(fmt.Sprintf("-%d", data.Cost.TotalLinesRemoved), Red))
		}
		parts = append(parts, strings.Join(changes, " "))
	}

	// 会话时长
	if cfg.IsLine1Enabled("duration") && data.Cost.TotalDurationMs > 0 {
		parts = append(parts, Colorize("⏱️ "+formatter.FormatDuration(data.Cost.TotalDurationMs), Blue))
	}

	// Token 统计
	if cfg.IsLine1Enabled("tokens") && (data.ContextWindow.TotalInputTokens > 0 || data.ContextWindow.TotalOutputTokens > 0) {
		in := formatter.FormatTokens(data.ContextWindow.TotalInputTokens)
		out := formatter.FormatTokens(data.ContextWindow.TotalOutputTokens)
		parts = append(parts, Colorize(fmt.Sprintf("📥%s 📤%s", in, out), Cyan))
	}

	// Nyan Cat 动画 + 处理状态指示器
	if cfg.IsLine1Enabled("nyan") {
		nyan := animation.NyanFrame()
		if indicator := processingIndicator(); indicator != "" {
			nyan += indicator
		}
		parts = append(parts, nyan)
	}

	// 心跳动画
	if cfg.IsLine1Enabled("heartbeat") {
		parts = append(parts, Colorize(animation.Heartbeat(), Red))
	}

	termWidth := GetTerminalWidth()
	return wrapParts(parts, sep, termWidth)
}

// renderLine2 渲染第二行: 统计信息
func renderLine2(sep string, cfg *config.Config) string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	binaryDir := filepath.Dir(execPath)

	info, err := stats.GetStatsInfo(binaryDir)
	if err != nil || info == nil {
		return ""
	}

	var parts []string

	if cfg.IsLine2Enabled("codingDays") && info.CodingDays > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("📅 %d天", info.CodingDays), Magenta))
	}
	if cfg.IsLine2Enabled("activeDays") && info.ActiveDays > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("🔥 %d天", info.ActiveDays), Green))
	}
	if cfg.IsLine2Enabled("streak") && info.Streak > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("⚡ %d连", info.Streak), Yellow))
	}
	if cfg.IsLine2Enabled("sessions") && info.TotalSessions > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("💬 %d会话", info.TotalSessions), Blue))
	}
	if cfg.IsLine2Enabled("messages") && info.TotalMessages > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("🗣️ %d消息", info.TotalMessages), Cyan))
	}
	if cfg.IsLine2Enabled("todayMessages") && info.TodayMessages > 0 {
		parts = append(parts, Colorize(fmt.Sprintf("📈 今日%d", info.TodayMessages), Cyan))
	}
	if cfg.IsLine2Enabled("peakHour") && info.HasPeakHour {
		emoji := peakHourEmoji(info.PeakHour)
		parts = append(parts, Colorize(fmt.Sprintf("%s %d点", emoji, info.PeakHour), Blue))
	}

	// 成就
	if cfg.IsLine2Enabled("achievement") {
		if achievement := getAchievement(info); achievement != "" {
			parts = append(parts, Colorize(achievement, Yellow))
		}
	}

	// 随机状态
	if cfg.IsLine2Enabled("randomStatus") {
		parts = append(parts, Colorize(animation.RandomStatus(), Cyan))
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, sep)
}

// processingIndicator 读取 hook 写入的状态文件, 返回处理状态指示器
// 处理中返回 "⏳", 处理完成返回 "⌛💯"
func processingIndicator() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	binaryDir := filepath.Dir(execPath)

	if state.IsProcessing(binaryDir) {
		return "⏳"
	}
	return "⌛💯"
}

// calcContextPercent 计算上下文使用百分比
func calcContextPercent(data *model.SessionData) float64 {
	if data.ContextWindow.ContextWindowSize <= 0 || data.ContextWindow.CurrentUsage == nil {
		return 0
	}
	usage := data.ContextWindow.CurrentUsage
	total := usage.InputTokens + usage.CacheCreationInputTokens + usage.CacheReadInputTokens
	return float64(total) / float64(data.ContextWindow.ContextWindowSize) * 100
}

// peakHourEmoji 根据小时返回时段 emoji
func peakHourEmoji(hour int) string {
	switch {
	case hour >= 22 || hour < 5:
		return "🌙"
	case hour >= 18:
		return "🌆"
	case hour >= 12:
		return "☀️"
	default:
		return "🌅"
	}
}

// getAchievement 根据统计数据返回成就徽章
func getAchievement(info *model.StatsInfo) string {
	switch {
	case info.TotalMessages >= 1000:
		return "🏆 千言万语"
	case info.TotalMessages >= 500:
		return "🥇 消息达人"
	case info.TotalMessages >= 100:
		return "🥈 话唠新星"
	}
	switch {
	case info.TotalSessions >= 100:
		return "👑 会话之王"
	case info.TotalSessions >= 50:
		return "⭐ 会话专家"
	}
	switch {
	case info.Streak >= 30:
		return "🔥 月度坚持"
	case info.Streak >= 7:
		return "💪 周度坚持"
	case info.Streak >= 3:
		return "✊ 三连击"
	}
	if info.ActiveDays >= 30 {
		return "🎖️ 老用户"
	}
	return ""
}
