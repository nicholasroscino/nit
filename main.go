package main

import (
	"nit/cli"
	"nit/utils"
	"os"
)

func main() {
	path, err := os.Getwd()
	utils.Check(err, "Unable to get current working directory\n")

	//err = os.Setenv("NIT_DEBUG", "1")
	utils.Check(err, "Unable to set NIT_DEBUG environment variable\n")

	commandDispatcher := cli.NewCliCommandDispatcher()
	commandDispatcher.Init()
	commandDispatcher.DispatchCommand(path, os.Args)
}
