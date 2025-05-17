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

	//val, err := commands.CatFileCommand(nitFolder, "3ccc5503fe7775a47905c5bf7be999189234d9c8")
	//
	//if err != nil {
	//	log.Fatal(err.Error())
	//
	//}

	//log.Println(val)
	err = commands.AddCommand(nitFolder, "/main.go")

	// Check if the file was added to the index
	if err != nil {
		log.Fatal(err.Error())
	}
}
