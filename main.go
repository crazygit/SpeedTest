package main

import (
	"github.com/crazygit/simple-download-tool/cmd"
	"github.com/crazygit/simple-download-tool/util"
	"os"
)

func main() {
	util.Log.WithField("args", os.Args[1:]).Info("App Started")
	cmd.Execute()
}
