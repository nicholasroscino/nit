package main

import (
	"log"
	"nit/commands"
	"nit/utils"
	"os"
)

func isNitRepo(path string) bool {
	if _, err := os.Stat(path + "/.nit"); os.IsNotExist(err) {
		return false
	}

	return true
}

func main() {
	rootFolder, err := os.Getwd()
	utils.Check(err, "Unable to get current working directory\n")

	if !isNitRepo(rootFolder) {
		log.Fatal("Not a nit repository")
	}

	commands.HashObjectCommand(rootFolder, "file.txt")
}
