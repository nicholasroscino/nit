package commit

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func commitCliCommand(cmd string) (*flag.FlagSet, *string, *string) {
	commitFlag := flag.NewFlagSet(cmd, flag.ExitOnError)

	message := commitFlag.String("m", "", "Commit message")
	author := commitFlag.String("a", "", "Author of the commit")

	return commitFlag, message, author
}

func handleCommitCommand(commitFlag *flag.FlagSet, message *string, author *string, projectPath string, args []string) {
	err := commitFlag.Parse(args)
	utils.Check(err, "Error parsing commit command arguments")
	if *message == "" || *author == "" {
		log.Fatal("commit command requires a commit message and an author.\n")
	}
	_, err = commitCommand(projectPath, *author, *message)
	utils.Check(err, "Error executing commit command")
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "commit"
	commitFlag, message, author := commitCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleCommitCommand(commitFlag, message, author, projectPath, args)
		},
		Name: command,
	}
}
