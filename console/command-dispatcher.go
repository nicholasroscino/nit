package console

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func hashCommand() (*flag.FlagSet, *string) {
	fooCmd := flag.NewFlagSet("hash", flag.ExitOnError)
	path := fooCmd.String("path", "", "Path to the file to hash")

	return fooCmd, path
}
func catCommand() (*flag.FlagSet, *string) {
	fooCmd := flag.NewFlagSet("init", flag.ExitOnError)

	hash := fooCmd.String("hash", "", "hash to display")

	return fooCmd, hash
}
func commitCommand() (*flag.FlagSet, *string, *string) {
	fooCmd := flag.NewFlagSet("init", flag.ExitOnError)

	message := fooCmd.String("m", "", "Commit message")
	author := fooCmd.String("a", "", "Author of the commit")

	return fooCmd, message, author
}

func addCommand() *flag.FlagSet {
	fooCmd := flag.NewFlagSet("init", flag.ExitOnError)

	return fooCmd
}

func initCommand() *flag.FlagSet {
	fooCmd := flag.NewFlagSet("init", flag.ExitOnError)
	return fooCmd
}

func DispatchCommand(projectPath string, osArgs []string) {
	initFlag := initCommand()
	catFlag, hashValue := catCommand()
	commitFlag, message, author := commitCommand()
	addFLag := addCommand()
	hashFlag, fullPathToHash := hashCommand()

	switch osArgs[1] {
	case "init":
		err := initFlag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing init command arguments")
		commands.InitCommand(projectPath)
	case "cat":
		err := catFlag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing cat command arguments")
		if *hashValue == "" {
			log.Fatal("cat command requires a hash value.\n")
		}
		value, err := commands.CatFileCommand(projectPath, *hashValue)
		utils.Check(err, "Error executing cat command")
		println(value)
	case "hash-object":
		err := hashFlag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing hash-object command arguments")
		if *fullPathToHash == "" {
			log.Fatal("hash-object command requires a file path.\n")
		}
		commands.HashObjectCommand(projectPath, *fullPathToHash)
	case "write-tree":
		err := initFlag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing write-tree command arguments")
		log.Println(commands.WriteTreeCommand(projectPath))
	case "add":
		err := addFLag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing add command arguments")
		files := addFLag.Args()

		if len(files) == 0 {
			log.Fatal("add command requires a file path.\n")
		}

		for _, pathToAdd := range files {
			err = commands.AddCommand(projectPath, pathToAdd)
			utils.Check(err, "Error executing add command")
		}

	case "commit":
		err := commitFlag.Parse(osArgs[2:])
		utils.Check(err, "Error parsing commit command arguments")
		if *message == "" || *author == "" {
			log.Fatal("commit command requires a commit message and an author.\n")
		}
		_, err = commands.CommitCommand(projectPath, *author, *message)
		utils.Check(err, "Error executing commit command")
	default:
		log.Fatal("Unknown command: " + osArgs[1] + "\n")
	}
}
