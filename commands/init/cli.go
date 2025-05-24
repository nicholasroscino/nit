package init

import (
	"flag"
	"nit/commands"
	"nit/utils"
)

func initCliCommand(cmd string) *flag.FlagSet {
	initFlag := flag.NewFlagSet(cmd, flag.ExitOnError)
	return initFlag
}

func handleInitCommand(initFlag *flag.FlagSet, projectPath string, args []string) {
	err := initFlag.Parse(args)
	utils.Check(err, "Error parsing init command arguments")
	initCommand(projectPath)
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "init"
	initFlag := initCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			if len(args) > 0 {
				utils.Check(nil, "init command does not require any arguments")
			}
			handleInitCommand(initFlag, projectPath, args)
		},
		Name: command,
	}
}
