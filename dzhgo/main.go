package main

import (
	"github.com/gzdzh-cn/dzhcore/dzhgo/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()
	cmd.Root.Run(ctx)
}
