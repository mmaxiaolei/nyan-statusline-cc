package animation

import (
	"fmt"
	"strings"
	"testing"
)

// TestRainbowProgressBar_ZeroPercent 验证 0% 时全部为空槽
func TestRainbowProgressBar_ZeroPercent(t *testing.T) {
	bar := RainbowProgressBar(0, 10)
	if bar == "" {
		t.Error("RainbowProgressBar(0, 10) should not return empty string")
	}
	// 0% 时不应包含填充字符 "█"
	if strings.Contains(bar, "█") {
		t.Error("0% bar should not contain filled blocks")
	}
	// 应包含空槽字符 "░"
	if !strings.Contains(bar, "░") {
		t.Error("0% bar should contain empty slots")
	}
}

// TestRainbowProgressBar_FullPercent 验证 100% 时全部为填充块
func TestRainbowProgressBar_FullPercent(t *testing.T) {
	bar := RainbowProgressBar(100, 10)
	if !strings.Contains(bar, "█") {
		t.Error("100% bar should contain filled blocks")
	}
	if strings.Contains(bar, "░") {
		t.Error("100% bar should not contain empty slots")
	}
}

// TestRainbowProgressBar_Width 验证不同宽度的进度条都能正常生成
func TestRainbowProgressBar_Width(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		width   int
	}{
		{"width=5 half", 50, 5},
		{"width=20 full", 100, 20},
		{"width=1 full", 100, 1},
		{"zero width fallback", 50, 0},    // 应回退到默认宽度 10
		{"negative width fallback", 50, -1}, // 应回退到默认宽度 10
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := RainbowProgressBar(tt.percent, tt.width)
			if bar == "" {
				t.Error("progress bar should not be empty")
			}
			// 所有进度条都应以 Reset 结尾
			if !strings.HasSuffix(bar, "\033[0m") {
				t.Error("progress bar should end with ANSI reset")
			}
		})
	}
}

// TestRainbowProgressBar_ContainsRainbowColors 验证 100% 进度条包含彩虹色 ANSI 256 色序列
func TestRainbowProgressBar_ContainsRainbowColors(t *testing.T) {
	bar := RainbowProgressBar(100, 14)
	for _, code := range rainbow256 {
		seq := fmt.Sprintf("\033[38;5;%dm", code)
		if !strings.Contains(bar, seq) {
			t.Errorf("100%% bar (width=14) should contain color code %d", code)
		}
	}
}

// TestRainbowProgressBar_HalfFilled 验证 50% 时同时包含填充块和空槽
func TestRainbowProgressBar_HalfFilled(t *testing.T) {
	bar := RainbowProgressBar(50, 10)
	if !strings.Contains(bar, "█") {
		t.Error("50% bar should contain filled blocks")
	}
	if !strings.Contains(bar, "░") {
		t.Error("50% bar should contain empty slots")
	}
}

// TestHeartbeat_NotEmpty 验证心跳动画输出非空
func TestHeartbeat_NotEmpty(t *testing.T) {
	hb := Heartbeat()
	if hb == "" {
		t.Error("Heartbeat() should not return empty string")
	}
}

// TestHeartbeat_ValidFrame 验证心跳返回的是有效帧
func TestHeartbeat_ValidFrame(t *testing.T) {
	hb := Heartbeat()
	valid := false
	for _, frame := range heartbeatFrames {
		if hb == frame {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Heartbeat() returned unexpected value: %q", hb)
	}
}

// TestRandomStatus_NotEmpty 验证随机状态输出非空
func TestRandomStatus_NotEmpty(t *testing.T) {
	status := RandomStatus()
	if status == "" {
		t.Error("RandomStatus() should not return empty string")
	}
}

// TestRandomStatus_ValidStatus 验证随机状态返回的是有效状态文字
func TestRandomStatus_ValidStatus(t *testing.T) {
	status := RandomStatus()
	valid := false
	for _, s := range statusMessages {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("RandomStatus() returned unexpected value: %q", status)
	}
}
