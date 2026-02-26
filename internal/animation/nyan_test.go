package animation

import (
	"fmt"
	"strings"
	"testing"
)

// ansi256Color 返回 ANSI 256 色前景色序列
func ansi256Color(code int) string {
	return fmt.Sprintf("\033[38;5;%dm", code)
}

// TestRainbow256_SevenColors 验证彩虹色方案为 7 色 (与 NyanProgressBar 一致)
func TestRainbow256_SevenColors(t *testing.T) {
	if len(rainbow256) != 7 {
		t.Errorf("rainbow256 should have 7 colors, got %d", len(rainbow256))
	}
}

// TestNyanFrame_NotEmpty 验证 NyanFrame 输出非空
func TestNyanFrame_NotEmpty(t *testing.T) {
	frame := NyanFrame()
	if frame == "" {
		t.Error("NyanFrame() should not return empty string")
	}
}

// TestNyanFrame_ContainsCatEmoji 验证输出包含猫咪 emoji
func TestNyanFrame_ContainsCatEmoji(t *testing.T) {
	// 遍历所有帧, 确保每帧都包含猫咪 emoji
	for i := range len(catFrames) {
		frame := nyanFrameAt(i)
		hasCat := false
		for _, cat := range catFrames {
			if strings.Contains(frame, cat) {
				hasCat = true
				break
			}
		}
		if !hasCat {
			t.Errorf("frame %d should contain a cat emoji, got: %q", i, frame)
		}
	}
}

// TestNyanFrameAt_ColorShift 验证不同帧的彩虹尾巴颜色偏移
// 帧 0 的第一个色块应为 rainbow256[0], 帧 1 的第一个色块应为 rainbow256[1]
func TestNyanFrameAt_ColorShift(t *testing.T) {
	for frameIdx := range len(rainbow256) {
		frame := nyanFrameAt(frameIdx)
		expectedFirst := ansi256Color(rainbow256[frameIdx%len(rainbow256)])
		if !strings.HasPrefix(frame, expectedFirst) {
			t.Errorf("frame %d: expected tail to start with color %d, got prefix: %q",
				frameIdx, rainbow256[frameIdx%len(rainbow256)], frame[:min(20, len(frame))])
		}
	}
}

// TestNyanFrameAt_TailContainsAllColors 验证每帧尾巴包含全部 7 种颜色
func TestNyanFrameAt_TailContainsAllColors(t *testing.T) {
	frame := nyanFrameAt(0)
	for _, code := range rainbow256 {
		colorSeq := ansi256Color(code)
		if !strings.Contains(frame, colorSeq) {
			t.Errorf("frame should contain ANSI 256 color %d, got: %q", code, frame)
		}
	}
}

// TestNyanFrameAt_ContainsReset 验证尾巴末尾包含 ANSI Reset
func TestNyanFrameAt_ContainsReset(t *testing.T) {
	frame := nyanFrameAt(0)
	if !strings.Contains(frame, "\033[0m") {
		t.Error("frame should contain ANSI reset sequence")
	}
}

// TestNyanFrameAt_ContainsStar 验证每帧包含星星点缀
func TestNyanFrameAt_ContainsStar(t *testing.T) {
	for i := range len(starFrames) {
		frame := nyanFrameAt(i)
		expected := starFrames[i%len(starFrames)]
		if !strings.Contains(frame, expected) {
			t.Errorf("frame %d should contain star %q", i, expected)
		}
	}
}

// TestNyanFrameAt_FrameCycle 验证帧循环: 相同帧索引产生相同输出
func TestNyanFrameAt_FrameCycle(t *testing.T) {
	for i := range len(catFrames) {
		first := nyanFrameAt(i)
		second := nyanFrameAt(i)
		if first != second {
			t.Errorf("same frameIdx %d should produce identical output", i)
		}
	}
}
