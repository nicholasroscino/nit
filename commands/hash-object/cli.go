package hash_object

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func hashCliCommand(cmd string) (*flag.FlagSet, *string) {
	hashFlag := flag.NewFlagSet(cmd, flag.ExitOnError)
	path := hashFlag.String("path", "", "Path to the file to hash")

	return hashFlag, path
}

func handleHashCommand(hashFlag *flag.FlagSet, fullPathToHash *string, projectPath string, args []string) {
	err := hashFlag.Parse(args)
	utils.Check(err, "Error parsing hash-object command arguments")
	if *fullPathToHash == "" {
		log.Fatal("hash-object command requires a file path.\n")
	}
	hashObjectCommand(projectPath, *fullPathToHash)
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "hash"
	hashFlag, fullPathToHash := hashCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleHashCommand(hashFlag, fullPathToHash, projectPath, args)
		},
		Name: command,
	}
}
