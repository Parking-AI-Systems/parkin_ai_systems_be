package main

import (
	_ "parkin-ai-system/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"parkin-ai-system/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
