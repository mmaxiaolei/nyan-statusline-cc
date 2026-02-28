package config

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// menuItem 菜单项
type menuItem struct {
	label   string
	key     string
	line    int  // 1 或 2, 标识属于哪一行; 0 表示 line2Enabled 开关
	enabled bool
	header  bool
}

// RunInteractive 启动交互式配置界面
func RunInteractive(dir string) error {
	cfg := Load(dir)
	items := buildMenuItems(cfg)
	cursor := nextSelectable(items, 0, 1)

	oldState, err := enableRawMode()
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer disableRawMode(oldState)

	// 预打印空行, 为首次 renderMenu 的光标上移腾出空间
	totalLines := len(items) + 4
	for i := 0; i < totalLines; i++ {
		fmt.Println()
	}

	saved := false
	for {
		renderMenu(items, cursor)
		switch readKey() {
		case "up":
			cursor = nextSelectable(items, cursor, -1)
		case "down":
			cursor = nextSelectable(items, cursor, 1)
		case "toggle":
			if !items[cursor].header {
				items[cursor].enabled = !items[cursor].enabled
			}
		case "save":
			applyToConfig(cfg, items)
			if err := Save(dir, cfg); err != nil {
				return err
			}
			saved = true
			fallthrough
		case "quit":
			// 清除菜单区域
			totalLines := len(items) + 4 // 标题 + items + 底部提示 + 空行
			fmt.Printf("\033[%dA\033[J", totalLines)
			if saved {
				fmt.Println("🐱 配置已保存 meow~")
			} else {
				fmt.Println("🐱 已取消 meow~")
			}
			return nil
		}
	}
}

func buildMenuItems(cfg *Config) []menuItem {
	var items []menuItem
	items = append(items, menuItem{label: "── Line 1 字段 ──", header: true})
	for _, f := range Line1Fields {
		items = append(items, menuItem{
			label: f.Label, key: f.Key, line: 1,
			enabled: cfg.IsLine1Enabled(f.Key),
		})
	}
	items = append(items, menuItem{label: "── Line 2 ──", header: true})
	items = append(items, menuItem{
		label: "✨ 启用第二行", key: "line2_enabled",
		enabled: cfg.Line2Enabled,
	})
	for _, f := range Line2Fields {
		items = append(items, menuItem{
			label: f.Label, key: f.Key, line: 2,
			enabled: cfg.IsLine2Enabled(f.Key),
		})
	}
	return items
}

func applyToConfig(cfg *Config, items []menuItem) {
	for _, it := range items {
		if it.header {
			continue
		}
		switch {
		case it.key == "line2_enabled":
			cfg.Line2Enabled = it.enabled
		case it.line == 1:
			cfg.Line1[it.key] = it.enabled
		case it.line == 2:
			cfg.Line2[it.key] = it.enabled
		}
	}
}

func nextSelectable(items []menuItem, cur, dir int) int {
	n := len(items)
	i := cur + dir
	for i >= 0 && i < n {
		if !items[i].header {
			return i
		}
		i += dir
	}
	return cur
}

func renderMenu(items []menuItem, cursor int) {
	// 移动光标到菜单起始位置并清除
	totalLines := len(items) + 4
	fmt.Printf("\033[%dA\033[J", totalLines)

	fmt.Println("\033[95m\033[1m🐱 Nyan Statusline 配置 meow~\033[0m")
	fmt.Println()
	for i, it := range items {
		prefix := "  "
		if i == cursor {
			prefix = "\033[96m> \033[0m"
		}
		if it.header {
			fmt.Printf("  \033[90m%s\033[0m\n", it.label)
			continue
		}
		check := "\033[92m✅\033[0m"
		if !it.enabled {
			check = "\033[90m⬜\033[0m"
		}
		fmt.Printf("%s%s %s\n", prefix, check, it.label)
	}
	fmt.Println()
	fmt.Println("\033[90m↑↓ 移动  空格 切换  Enter 保存  q 取消\033[0m")
}

// --- 终端 raw mode (syscall, 无外部依赖) ---

type termState struct {
	termios syscall.Termios
}

func enableRawMode() (*termState, error) {
	var orig syscall.Termios
	if err := ioctl(syscall.Stdin, syscall.TIOCGETA, &orig); err != nil {
		return nil, err
	}
	raw := orig
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if err := ioctl(syscall.Stdin, syscall.TIOCSETA, &raw); err != nil {
		return nil, err
	}
	return &termState{termios: orig}, nil
}

func disableRawMode(state *termState) {
	_ = ioctl(syscall.Stdin, syscall.TIOCSETA, &state.termios)
}

func ioctl(fd int, req uint, arg *syscall.Termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(fd), uintptr(req),
		uintptr(unsafe.Pointer(arg)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func readKey() string {
	buf := make([]byte, 3)
	n, _ := os.Stdin.Read(buf)
	if n == 0 {
		return ""
	}
	// ESC 序列 (方向键)
	if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return "up"
		case 'B':
			return "down"
		}
	}
	switch buf[0] {
	case ' ':
		return "toggle"
	case '\r', '\n':
		return "save"
	case 'q', 0x1b: // q 或 ESC
		return "quit"
	}
	return ""
}
