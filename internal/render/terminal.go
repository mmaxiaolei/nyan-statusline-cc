package render

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// winsize 对应 TIOCGWINSZ 返回的终端尺寸结构
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

const defaultTermWidth = 200

// GetTerminalWidth 获取当前终端列数.
// 优先通过 /dev/tty 发起 ioctl 查询 (即使 stdout/stdin 被 Claude Code 管道重定向也有效),
// 其次读取 COLUMNS 环境变量, 均失败时返回默认值 200.
func GetTerminalWidth() int {
	// 通过 /dev/tty 获取终端尺寸, 不依赖 stdout/stdin 是否为 tty
	tty, err := os.Open("/dev/tty")
	if err == nil {
		defer tty.Close()
		var ws winsize
		_, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			tty.Fd(),
			syscall.TIOCGWINSZ,
			uintptr(unsafe.Pointer(&ws)),
		)
		if errno == 0 && ws.Col > 0 {
			return int(ws.Col)
		}
	}
	// 备选: $COLUMNS 环境变量 (某些 shell 会自动设置)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}
	return defaultTermWidth
}
