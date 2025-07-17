package cmd

import (
	"github.com/gogf/gf/v2/os/gcmd"
)

// Root 根命令
var (
	Root = gcmd.Command{
		Name:  "dzhgo-cli",
		Usage: "dzhgo-cli [COMMAND] [OPTION]",
		Brief: "DzhGO 代码生成工具",
		Description: `
DzhGO 代码生成工具，用于快速生成控制器、模型和服务代码。

使用示例：
1. 初始化项目: dzhgo-cli init -n myproject
2. 生成代码: dzhgo-cli gen -a user -n user -m admin
`,
		Additional: `
运行 'dzhgo-cli COMMAND -h' 获取更多命令帮助信息。
`,
	}
)

func init() {
	// 添加所有命令到根命令
	Root.AddCommand(GenCode)
	Root.AddCommand(InitProject)
}
