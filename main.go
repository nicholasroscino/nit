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

	//commands.HashObjectCommand(nitFolder, rootFolder+"/main.go")
	ret, err := commands.CatFileCommand(nitFolder, "ac54cd8cfa2b1487e4b786ad29eced9829a2bd0b")

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(ret)
}
