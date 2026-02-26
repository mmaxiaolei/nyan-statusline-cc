package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nyan-statusline-cc/internal/parser"
	"github.com/nyan-statusline-cc/internal/render"
	"github.com/nyan-statusline-cc/internal/state"
)

func main() {
	// --state processing/completed: hooks 调用模式, 写入状态后退出
	if len(os.Args) == 3 && os.Args[1] == "--state" {
		binaryDir := filepath.Dir(os.Args[0])
		if err := state.SetStatus(binaryDir, os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "set state error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 默认模式: 从 stdin 读取会话数据并渲染状态栏
	data, err := parser.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(render.Render(data))
}
