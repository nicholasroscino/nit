package write_tree

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func writeTreeCliCommand(cmd string) *flag.FlagSet {
	writeTreeFlag := flag.NewFlagSet(cmd, flag.ExitOnError)
	return writeTreeFlag
}

func handleWriteTreeCommand(writeTreeFlag *flag.FlagSet, projectPath string, args []string) {
	err := writeTreeFlag.Parse(args)
	utils.Check(err, "Error parsing write-tree command arguments")
	log.Println(WriteTreeCommand(projectPath))
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "write-tree"
	writeTreeFlag := writeTreeCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			if len(args) > 0 {
				utils.Check(nil, "write-tree command does not require any arguments")
			}

			handleWriteTreeCommand(writeTreeFlag, projectPath, args)
		},
		Name: command,
	}
}
