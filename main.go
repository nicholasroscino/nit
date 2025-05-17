package main

import (
	"log"
	"nit/commands"
	"nit/utils"
	"os"
)

func main() {
	rootFolder, err := os.Getwd()
	utils.Check(err, "Unable to get current working directory\n")

	nitFolder, err := utils.GetNitRepoFolder(rootFolder)

	if err != nil {
		log.Fatal(err.Error())
	}

	commands.HashObjectCommand(nitFolder, rootFolder+"/file.txt")
}
