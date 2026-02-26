package main

import (
	"fmt"
	"os"

	"github.com/nyan-statusline-cc/internal/parser"
	"github.com/nyan-statusline-cc/internal/render"
)

func main() {
	data, err := parser.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(render.Render(data))
}
