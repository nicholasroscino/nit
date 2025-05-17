package main

import (
	"nit/commands"
	"nit/utils"
	"os"
)

func main() {
	//println(os.Getenv("GOPATH"))
	val, err := os.Getwd()
	utils.Check(err, "Unable to get current working directory\n")

	commands.HashObjectCommand(val+"/main.go", val)
}
