package cat

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func catCliCommand(cmd string) (*flag.FlagSet, *string) {
	catCmd := flag.NewFlagSet(cmd, flag.ExitOnError)

	hash := catCmd.String("hash", "", "hash to display")

	return catCmd, hash
}

func handleCatCommand(catFlag *flag.FlagSet, hashValue *string, projectPath string, args []string) {
	err := catFlag.Parse(args)
	utils.Check(err, "Error parsing cat command arguments")
	if *hashValue == "" {
		log.Fatal("cat command requires a hash value.\n")
	}
	value, err := catFileCommand(projectPath, *hashValue)
	utils.Check(err, "Error executing cat command")
	println(value)
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "cat"
	catFlag, catParam := catCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleCatCommand(catFlag, catParam, projectPath, args)
		},
		Name: command,
	}
}
