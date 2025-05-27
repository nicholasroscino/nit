package log

import (
	"flag"
	"nit/commands"
	"nit/utils"
)

func logCliCommand(cmd string) *flag.FlagSet {
	logFlag := flag.NewFlagSet(cmd, flag.ExitOnError)

	return logFlag
}

func handleLogCommand(logFlag *flag.FlagSet, projectPath string, args []string) {
	err := logFlag.Parse(args)
	utils.Check(err, "Error parsing log command arguments")

	err = logCommand(projectPath)
	utils.Check(err, "Error executing log command")
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "log"
	logFlag := logCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleLogCommand(logFlag, projectPath, args)
		},
		Name: command,
	}
}
