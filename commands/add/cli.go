package add

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func addCliCommand(command string) *flag.FlagSet {
	addFlag := flag.NewFlagSet(command, flag.ExitOnError)

	return addFlag
}

func handleAddCommand(addFlag *flag.FlagSet, projectPath string, args []string) {
	err := addFlag.Parse(args)
	utils.Check(err, "Error parsing add command arguments")

	files := addFlag.Args()

	if len(files) == 0 {
		log.Fatal("add command requires a file path.\n")
	}

	for _, pathToAdd := range files {
		err = addCommand(projectPath, pathToAdd)
		utils.Check(err, "Error executing add command")
	}
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "add"

	addFlag := addCliCommand(command)
	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleAddCommand(addFlag, projectPath, args)
		},
		Name: command,
	}
}
