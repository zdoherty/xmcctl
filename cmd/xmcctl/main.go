package main

import (
	"git.poundadm.net/anachronism/xmcctl/cmd/xmcctl/cmds"
)

func main() {
	root := cmds.New()
	root.Execute()
}
