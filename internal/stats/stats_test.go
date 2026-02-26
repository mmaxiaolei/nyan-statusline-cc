package stats

import (
	"testing"
	"time"

	"github.com/nyan-statusline-cc/internal/model"
)

// 固定测试时间: 2026-02-26
var testNow = time.Date(2026, 2, 26, 10, 0, 0, 0, time.UTC)

func TestCalcStreak_Empty(t *testing.T) {
	got := calcStreak(nil, testNow)
	if got != 0 {
		t.Errorf("calcStreak(nil) = %d, want 0", got)
	}
}

func TestCalcStreak_OnlyToday(t *testing.T) {
	activity := []model.DailyActivity{
		{Date: "2026-02-26", MessageCount: 5},
	}
	got := calcStreak(activity, testNow)
	if got != 1 {
		t.Errorf("calcStreak(today only) = %d, want 1", got)
	}
}

func TestCalcStreak_OnlyYesterday(t *testing.T) {
	activity := []model.DailyActivity{
		{Date: "2026-02-25", MessageCount: 3},
	}
	got := calcStreak(activity, testNow)
	if got != 1 {
		t.Errorf("calcStreak(yesterday only) = %d, want 1", got)
	}
}

func TestCalcStreak_ConsecutiveDays(t *testing.T) {
	// 连续 5 天: 2/22 ~ 2/26
	activity := []model.DailyActivity{
		{Date: "2026-02-22"},
		{Date: "2026-02-23"},
		{Date: "2026-02-24"},
		{Date: "2026-02-25"},
		{Date: "2026-02-26"},
	}
	got := calcStreak(activity, testNow)
	if got != 5 {
		t.Errorf("calcStreak(5 consecutive) = %d, want 5", got)
	}
}

func TestCalcStreak_BrokenStreak(t *testing.T) {
	// 2/26, 2/25, 2/23 (2/24 缺失, streak 应为 2)
	activity := []model.DailyActivity{
		{Date: "2026-02-23"},
		{Date: "2026-02-25"},
		{Date: "2026-02-26"},
	}
	got := calcStreak(activity, testNow)
	if got != 2 {
		t.Errorf("calcStreak(broken) = %d, want 2", got)
	}
}

func TestCalcStreak_OldDatesOnly(t *testing.T) {
	// 最近日期是 3 天前, streak 应为 0
	activity := []model.DailyActivity{
		{Date: "2026-02-20"},
		{Date: "2026-02-21"},
		{Date: "2026-02-22"},
	}
	got := calcStreak(activity, testNow)
	if got != 0 {
		t.Errorf("calcStreak(old dates) = %d, want 0", got)
	}
}

func TestCalcStreak_DuplicateDates(t *testing.T) {
	// 重复日期不应影响 streak 计算
	activity := []model.DailyActivity{
		{Date: "2026-02-25"},
		{Date: "2026-02-25"},
		{Date: "2026-02-26"},
		{Date: "2026-02-26"},
	}
	got := calcStreak(activity, testNow)
	if got != 2 {
		t.Errorf("calcStreak(duplicates) = %d, want 2", got)
	}
}

func TestComputeStatsInfo_NilCache(t *testing.T) {
	info := ComputeStatsInfo(nil, testNow)
	if info == nil {
		t.Fatal("ComputeStatsInfo(nil) should return non-nil empty StatsInfo")
	}
	if info.TotalSessions != 0 || info.TotalMessages != 0 {
		t.Error("nil cache should produce zero-value StatsInfo")
	}
}

func TestComputeStatsInfo_EmptyCache(t *testing.T) {
	cache := &model.StatsCache{}
	info := ComputeStatsInfo(cache, testNow)
	if info.ActiveDays != 0 {
		t.Errorf("ActiveDays = %d, want 0", info.ActiveDays)
	}
	if info.Streak != 0 {
		t.Errorf("Streak = %d, want 0", info.Streak)
	}
	if info.HasPeakHour {
		t.Error("HasPeakHour should be false for empty cache")
	}
}

func TestComputeStatsInfo_FullData(t *testing.T) {
	cache := &model.StatsCache{
		FirstSessionDate: "2026-02-20T08:00:00Z",
		TotalSessions:    10,
		TotalMessages:    200,
		DailyActivity: []model.DailyActivity{
			{Date: "2026-02-20", MessageCount: 30, SessionCount: 2},
			{Date: "2026-02-21", MessageCount: 40, SessionCount: 3},
			{Date: "2026-02-25", MessageCount: 50, SessionCount: 2},
			{Date: "2026-02-26", MessageCount: 80, SessionCount: 3},
		},
		HourCounts: map[string]int{"9": 15, "14": 25, "22": 10},
	}

	info := ComputeStatsInfo(cache, testNow)

	if info.TotalSessions != 10 {
		t.Errorf("TotalSessions = %d, want 10", info.TotalSessions)
	}
	if info.TotalMessages != 200 {
		t.Errorf("TotalMessages = %d, want 200", info.TotalMessages)
	}
	if info.ActiveDays != 4 {
		t.Errorf("ActiveDays = %d, want 4", info.ActiveDays)
	}
	// 2/20 ~ 2/26 = 7 天
	if info.CodingDays != 7 {
		t.Errorf("CodingDays = %d, want 7", info.CodingDays)
	}
	// 连续: 2/26, 2/25 (2/24 缺失), streak = 2
	if info.Streak != 2 {
		t.Errorf("Streak = %d, want 2", info.Streak)
	}
}

func TestComputeStatsInfo_TodayStats(t *testing.T) {
	cache := &model.StatsCache{
		DailyActivity: []model.DailyActivity{
			{Date: "2026-02-26", MessageCount: 42, SessionCount: 7},
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if info.TodayMessages != 42 {
		t.Errorf("TodayMessages = %d, want 42", info.TodayMessages)
	}
	if info.TodaySessions != 7 {
		t.Errorf("TodaySessions = %d, want 7", info.TodaySessions)
	}
}

func TestComputeStatsInfo_TodayNotPresent(t *testing.T) {
	cache := &model.StatsCache{
		DailyActivity: []model.DailyActivity{
			{Date: "2026-02-20", MessageCount: 10, SessionCount: 1},
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if info.TodayMessages != 0 {
		t.Errorf("TodayMessages = %d, want 0", info.TodayMessages)
	}
	if info.TodaySessions != 0 {
		t.Errorf("TodaySessions = %d, want 0", info.TodaySessions)
	}
}

func TestComputeStatsInfo_PeakHour(t *testing.T) {
	cache := &model.StatsCache{
		HourCounts: map[string]int{
			"9":  15,
			"14": 25,
			"22": 10,
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if !info.HasPeakHour {
		t.Fatal("HasPeakHour should be true")
	}
	if info.PeakHour != 14 {
		t.Errorf("PeakHour = %d, want 14", info.PeakHour)
	}
	if info.PeakCount != 25 {
		t.Errorf("PeakCount = %d, want 25", info.PeakCount)
	}
}

func TestComputeStatsInfo_PeakHourTieBreak(t *testing.T) {
	// count 相同时应取较小 hour
	cache := &model.StatsCache{
		HourCounts: map[string]int{
			"10": 20,
			"15": 20,
			"3":  20,
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if !info.HasPeakHour {
		t.Fatal("HasPeakHour should be true")
	}
	if info.PeakHour != 3 {
		t.Errorf("PeakHour = %d, want 3 (smallest hour on tie)", info.PeakHour)
	}
}

func TestComputeStatsInfo_PeakHourInvalidKey(t *testing.T) {
	// 非法 key 应被跳过
	cache := &model.StatsCache{
		HourCounts: map[string]int{
			"abc": 100,
			"9":   5,
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if !info.HasPeakHour {
		t.Fatal("HasPeakHour should be true")
	}
	if info.PeakHour != 9 {
		t.Errorf("PeakHour = %d, want 9", info.PeakHour)
	}
	if info.PeakCount != 5 {
		t.Errorf("PeakCount = %d, want 5", info.PeakCount)
	}
}

func TestComputeStatsInfo_PeakHourAllInvalidKeys(t *testing.T) {
	cache := &model.StatsCache{
		HourCounts: map[string]int{
			"abc": 100,
			"xyz": 50,
		},
	}
	info := ComputeStatsInfo(cache, testNow)

	if info.HasPeakHour {
		t.Error("HasPeakHour should be false when all keys are invalid")
	}
}

func TestComputeStatsInfo_CodingDays_SameDay(t *testing.T) {
	// 首次使用就是今天, CodingDays 应为 1
	cache := &model.StatsCache{
		FirstSessionDate: "2026-02-26T08:00:00Z",
	}
	info := ComputeStatsInfo(cache, testNow)

	if info.CodingDays != 1 {
		t.Errorf("CodingDays = %d, want 1", info.CodingDays)
	}
}

func TestComputeStatsInfo_CodingDays_MultiDay(t *testing.T) {
	// 首次使用是 10 天前
	cache := &model.StatsCache{
		FirstSessionDate: "2026-02-16T20:00:00Z",
	}
	info := ComputeStatsInfo(cache, testNow)

	// 2/16 ~ 2/26 = 11 天 (含首尾)
	if info.CodingDays != 11 {
		t.Errorf("CodingDays = %d, want 11", info.CodingDays)
	}
}

func TestGetAchievement_Nil(t *testing.T) {
	got := GetAchievement(nil)
	if got != "" {
		t.Errorf("GetAchievement(nil) = %q, want empty", got)
	}
}

func TestGetAchievement_NoAchievement(t *testing.T) {
	info := &model.StatsInfo{
		TotalMessages: 10,
		TotalSessions: 5,
		Streak:        1,
		ActiveDays:    3,
	}
	got := GetAchievement(info)
	if got != "" {
		t.Errorf("GetAchievement = %q, want empty", got)
	}
}

func TestGetAchievement_MessageTiers(t *testing.T) {
	cases := []struct {
		name     string
		messages int
		want     string
	}{
		{"below_100", 99, ""},
		{"at_100", 100, "🥈 话唠新星"},
		{"at_500", 500, "🥇 消息达人"},
		{"at_1000", 1000, "🏆 千言万语"},
		{"above_1000", 5000, "🏆 千言万语"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info := &model.StatsInfo{TotalMessages: tc.messages}
			got := GetAchievement(info)
			if got != tc.want {
				t.Errorf("GetAchievement(messages=%d) = %q, want %q", tc.messages, got, tc.want)
			}
		})
	}
}

func TestGetAchievement_SessionTiers(t *testing.T) {
	cases := []struct {
		name     string
		sessions int
		want     string
	}{
		{"below_50", 49, ""},
		{"at_50", 50, "⭐ 会话专家"},
		{"at_100", 100, "👑 会话之王"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info := &model.StatsInfo{TotalSessions: tc.sessions}
			got := GetAchievement(info)
			if got != tc.want {
				t.Errorf("GetAchievement(sessions=%d) = %q, want %q", tc.sessions, got, tc.want)
			}
		})
	}
}

func TestGetAchievement_StreakTiers(t *testing.T) {
	cases := []struct {
		name   string
		streak int
		want   string
	}{
		{"below_3", 2, ""},
		{"at_3", 3, "✊ 三连击"},
		{"at_7", 7, "💪 周度坚持"},
		{"at_30", 30, "🔥 月度坚持"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info := &model.StatsInfo{Streak: tc.streak}
			got := GetAchievement(info)
			if got != tc.want {
				t.Errorf("GetAchievement(streak=%d) = %q, want %q", tc.streak, got, tc.want)
			}
		})
	}
}

func TestGetAchievement_ActiveDays(t *testing.T) {
	info := &model.StatsInfo{ActiveDays: 30}
	got := GetAchievement(info)
	if got != "🎖️ 老用户" {
		t.Errorf("GetAchievement(activeDays=30) = %q, want %q", got, "🎖️ 老用户")
	}
}

func TestGetAchievement_Priority(t *testing.T) {
	// 消息成就优先级高于会话成就
	info := &model.StatsInfo{
		TotalMessages: 1000,
		TotalSessions: 100,
		Streak:        30,
		ActiveDays:    30,
	}
	got := GetAchievement(info)
	if got != "🏆 千言万语" {
		t.Errorf("GetAchievement(all high) = %q, want %q", got, "🏆 千言万语")
	}
}
