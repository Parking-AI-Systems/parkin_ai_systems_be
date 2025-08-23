package main

import (
	_ "parkin-ai-system/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"parkin-ai-system/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
