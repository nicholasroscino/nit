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

	err = os.Setenv("NIT_DEBUG", "1")
	utils.Check(err, "Unable to set NIT_DEBUG environment variable\n")

	//commands.InitCommand(rootFolder)

	nitFolder, err := utils.GetNitRepoFolder(rootFolder)
	utils.Check(err, "Unable to find .nit folder in the current directory\n")

	err = commands.AddCommand(nitFolder, "main.go")
	if err != nil {
		log.Fatal(err.Error())
	}
	//
	//commitHash, err := commands.CatFileCommand(nitFolder, "79c9af3e42bfe5bbf4da21327b5324a0ad31a6e3")
	commitHash, err := commands.CommitCommand(nitFolder, "nick!", "third commit")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Commit hash:\n", commitHash)
}
