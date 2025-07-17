package cmd

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func TestPrintModName(t *testing.T) {

	var (
		str           = `module dzhgo-cli`
		characterMask = "module "
		result        = gstr.TrimLeftStr(str, characterMask)
	)
	t.Log("result", result)

	modName := ""
	if modData := gfile.GetContents("go.mod"); modData != "" {
		for _, line := range gstr.Split(modData, "\n") {
			if gstr.HasPrefix(line, "module ") {
				modName = gstr.Trim(gstr.TrimLeft(line, "module"))
				t.Log("modName", modName)
				modName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
				t.Log("modName", modName)
				break
			}
		}
	}
	t.Log("module name:", modName)
}
