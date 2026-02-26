package animation

import (
	"fmt"
	"strings"
	"time"
)

// ANSI 256 色彩虹色值
var rainbow256 = []int{196, 208, 226, 46, 51, 21, 93}

// RainbowProgressBar 生成彩虹渐变进度条
// Parameters:
//   - percent: 百分比 (0-100)
//   - width: 进度条字符宽度
//
// Return:
//   - string: 带 ANSI 彩虹色的进度条字符串
func RainbowProgressBar(percent float64, width int) string {
	if width <= 0 {
		width = 10
	}
	filled := min(int(float64(width)*percent/100), width)

	var b strings.Builder
	for i := range filled {
		colorIdx := min(i*len(rainbow256)/width, len(rainbow256)-1)
		fmt.Fprintf(&b, "\033[38;5;%dm█", rainbow256[colorIdx])
	}
	for range width - filled {
		b.WriteString("\033[90m░")
	}
	b.WriteString("\033[0m")
	return b.String()
}

// heartbeatFrames 随机动画帧序列
var heartbeatFrames = []string{"👻", "👹", "💗", "🎃"}

// statusMessages 随机状态文字池
var statusMessages = []string{
	"🚀 火力全开", "💡 灵感爆发", "🎯 专注模式", "⚡ 效率拉满",
	"🔮 魔法编程", "🎮 游戏时间", "☕ 咖啡时间", "🌙 深夜肝码",
	"🌅 早起的鸟", "🦾 AI 附体", "🧠 脑洞大开", "✨ 代码如诗",
}

// Heartbeat 返回当前心跳动画帧
// Return:
//   - string: 心跳 emoji
func Heartbeat() string {
	idx := int(time.Now().UnixMilli()/333) % len(heartbeatFrames)
	return heartbeatFrames[idx]
}

// RandomStatus 返回随机状态文字, 每分钟更新一次
// Return:
//   - string: 状态文字
func RandomStatus() string {
	idx := int(time.Now().Unix()/60) % len(statusMessages)
	return statusMessages[idx]
}
