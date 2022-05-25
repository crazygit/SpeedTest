package main

import (
	"github.com/crazygit/simple-download-tool/cmd"
	"github.com/crazygit/simple-download-tool/util"
	"os"
)

func main() {
	util.Log.WithField("args", os.Args).Info("App Started")
	cmd.Execute()
}
