// Package stats 负责读取和解析 Claude Code 的统计缓存数据
package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/nyan-statusline-cc/internal/model"
)

// GetStatsInfo 读取 stats-cache.json 并解析为统计摘要
// Parameters:
//   - binaryDir: 二进制文件所在目录 (stats-cache.json 同目录)
//
// Return:
//   - *model.StatsInfo: 统计摘要, 文件不存在时返回 nil
//   - error: 读取或解析错误
func GetStatsInfo(binaryDir string) (*model.StatsInfo, error) {
	statsPath := filepath.Join(binaryDir, "stats-cache.json")
	data, err := os.ReadFile(statsPath)
	if err != nil {
		return nil, nil
	}

	var cache model.StatsCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return ComputeStatsInfo(&cache, time.Now()), nil
}

// ComputeStatsInfo 根据缓存数据和当前时间计算统计摘要
// Parameters:
//   - cache: 统计缓存数据
//   - now: 当前时间 (用于计算天数和连续活跃)
//
// Return:
//   - *model.StatsInfo: 统计摘要
func ComputeStatsInfo(cache *model.StatsCache, now time.Time) *model.StatsInfo {
	if cache == nil {
		return &model.StatsInfo{}
	}

	info := &model.StatsInfo{
		TotalSessions: cache.TotalSessions,
		TotalMessages: cache.TotalMessages,
		ActiveDays:    len(cache.DailyActivity),
	}

	// 计算使用天数 (首日也算一天, 所以 +1)
	if cache.FirstSessionDate != "" {
		first, err := time.Parse(time.RFC3339, cache.FirstSessionDate)
		if err == nil {
			firstDate := truncateToDate(first)
			nowDate := truncateToDate(now)
			days := int(nowDate.Sub(firstDate).Hours()/24) + 1
			if days < 1 {
				days = 1
			}
			info.CodingDays = days
		}
	}

	// 计算连续活跃天数
	info.Streak = calcStreak(cache.DailyActivity, now)

	// 今日统计
	today := now.Format("2006-01-02")
	for _, day := range cache.DailyActivity {
		if day.Date == today {
			info.TodayMessages = day.MessageCount
			info.TodaySessions = day.SessionCount
			break
		}
	}

	// 最活跃时段: count 相同时取较小 hour, 保证确定性
	calcPeakHour(cache.HourCounts, info)

	return info
}

// calcPeakHour 从小时计数中找出最活跃时段
func calcPeakHour(hourCounts map[string]int, info *model.StatsInfo) {
	if len(hourCounts) == 0 {
		return
	}

	maxCount := 0
	peakHour := -1

	for h, c := range hourCounts {
		hour, err := strconv.Atoi(h)
		if err != nil {
			continue
		}
		if c > maxCount || (c == maxCount && (peakHour < 0 || hour < peakHour)) {
			maxCount = c
			peakHour = hour
		}
	}

	if peakHour >= 0 && maxCount > 0 {
		info.PeakHour = peakHour
		info.PeakCount = maxCount
		info.HasPeakHour = true
	}
}

// calcStreak 计算从今天/昨天开始的连续活跃天数
func calcStreak(activity []model.DailyActivity, now time.Time) int {
	if len(activity) == 0 {
		return 0
	}

	// 收集并去重日期
	dateSet := make(map[string]struct{}, len(activity))
	for _, d := range activity {
		dateSet[d.Date] = struct{}{}
	}

	dates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	// 连续天数必须从今天或昨天开始
	if dates[0] != today && dates[0] != yesterday {
		return 0
	}

	streak := 1
	for i := 0; i < len(dates)-1; i++ {
		curr, errC := time.Parse("2006-01-02", dates[i])
		prev, errP := time.Parse("2006-01-02", dates[i+1])
		if errC != nil || errP != nil {
			break
		}
		// 用日期减一天比较, 避免浮点精度问题
		if curr.AddDate(0, 0, -1).Equal(prev) {
			streak++
		} else {
			break
		}
	}
	return streak
}

// truncateToDate 将时间截断到日期 (去掉时分秒)
func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// GetAchievement 根据统计数据返回成就徽章文本
// Parameters:
//   - info: 统计摘要
//
// Return:
//   - string: 成就徽章文本, 无成就时返回空字符串
func GetAchievement(info *model.StatsInfo) string {
	if info == nil {
		return ""
	}

	// 消息数成就 (优先级最高)
	switch {
	case info.TotalMessages >= 1000:
		return "🏆 千言万语"
	case info.TotalMessages >= 500:
		return "🥇 消息达人"
	case info.TotalMessages >= 100:
		return "🥈 话唠新星"
	}

	// 会话数成就
	switch {
	case info.TotalSessions >= 100:
		return "👑 会话之王"
	case info.TotalSessions >= 50:
		return "⭐ 会话专家"
	}

	// 连续活跃成就
	switch {
	case info.Streak >= 30:
		return "🔥 月度坚持"
	case info.Streak >= 7:
		return "💪 周度坚持"
	case info.Streak >= 3:
		return "✊ 三连击"
	}

	// 活跃天数成就
	if info.ActiveDays >= 30 {
		return "🎖️ 老用户"
	}

	return ""
}
