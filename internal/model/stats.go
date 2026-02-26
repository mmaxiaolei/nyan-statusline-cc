// Package model 定义统计缓存的数据模型
package model

// StatsCache 表示 stats-cache.json 的数据结构
type StatsCache struct {
	FirstSessionDate string          `json:"firstSessionDate"`
	TotalSessions    int             `json:"totalSessions"`
	TotalMessages    int             `json:"totalMessages"`
	DailyActivity    []DailyActivity `json:"dailyActivity"`
	HourCounts       map[string]int  `json:"hourCounts"`
}

// DailyActivity 每日活动记录
type DailyActivity struct {
	Date         string `json:"date"`
	MessageCount int    `json:"messageCount"`
	SessionCount int    `json:"sessionCount"`
}

// StatsInfo 解析后的统计摘要
type StatsInfo struct {
	CodingDays    int
	ActiveDays    int
	Streak        int
	TotalSessions int
	TotalMessages int
	TodayMessages int
	TodaySessions int
	PeakHour      int
	PeakCount     int
	HasPeakHour   bool
}
