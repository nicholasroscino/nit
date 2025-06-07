package walk

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func walkCliCommand(cmd string) (*flag.FlagSet, *string) {
	walkFlag := flag.NewFlagSet(cmd, flag.ExitOnError)

	commitHash := walkFlag.String("c", "", "Commit hash to start walking from")

	return walkFlag, commitHash
}

func handleLogCommand(walkFlag *flag.FlagSet, commitHash *string, projectPath string, args []string) {
	err := walkFlag.Parse(args)
	utils.Check(err, "Error parsing log command arguments")

	if commitHash == nil || len(*commitHash) != 40 {
		log.Fatal("walk command requires a valid commit hash to move to.")
	}

	err = walkCommand(projectPath, commitHash)
	utils.Check(err, "Error executing log command")

	log.Println("You are currently in detached HEAD state. Current HEAD position:", *commitHash)
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "walk"
	walkFlag, commitHash := walkCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleLogCommand(walkFlag, commitHash, projectPath, args)
		},
		Name: command,
	}
}
