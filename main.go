package main

import (
	"os"
	"speedtest/cmd"
	"speedtest/util"
)

func main() {
	util.Log.WithField("args", os.Args[1:]).Info("App Started")
	cmd.Execute()
}
