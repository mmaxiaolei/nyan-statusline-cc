// Package animation 实现状态栏中的各类动画效果
package animation

import (
	"fmt"
	"time"
)

// catFrames 猫咪帧序列, 交替使用不同 emoji 表示动感
var catFrames = []string{"🐱", "😺", "🐱", "😸"}

// starFrames 星星点缀帧序列
var starFrames = []string{"✨", "⭐", "✨"}

// NyanFrame 返回当前 Nyan Cat 动画帧
// 猫咪使用 emoji, 彩虹尾巴使用 ANSI 256 色 7 色方案 (与 NyanProgressBar 一致)
// 颜色每帧偏移一位, 产生流畅的滚动效果
//
// Return:
//   - string: 当前帧的字符串表示
func NyanFrame() string {
	frameIdx := int(time.Now().UnixMilli()/250) % len(catFrames)
	return nyanFrameAt(frameIdx)
}

// nyanFrameAt 根据帧索引生成 Nyan Cat 动画帧
// Parameters:
//   - frameIdx: 帧索引, 同时用于驱动尾巴颜色偏移、猫咪和星星的切换
//
// Return:
//   - string: 当前帧的字符串表示
func nyanFrameAt(frameIdx int) string {
	// 彩虹尾巴: 使用 ANSI 256 色渲染 "█" 字符
	// 7 色方案: Red(196), Orange(208), Yellow(226), Green(46), Cyan(51), Blue(21), Violet(93)
	// 颜色每帧偏移一位, 产生滚动效果
	n := len(rainbow256)
	offset := frameIdx % n
	tail := make([]byte, 0, n*16) // 预分配足够空间
	for i := range n {
		idx := (i + offset) % n
		tail = fmt.Appendf(tail, "\033[38;5;%dm█", rainbow256[idx])
	}
	tail = append(tail, "\033[0m"...)

	cat := catFrames[frameIdx%len(catFrames)]
	star := starFrames[frameIdx%len(starFrames)]

	return string(tail) + cat + star
}
