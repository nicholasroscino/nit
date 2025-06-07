package walk

import (
	"flag"
	"log"
	"nit/commands"
	"nit/utils"
)

func walkCliCommand(cmd string) (*flag.FlagSet, *string, *string) {
	walkFlag := flag.NewFlagSet(cmd, flag.ExitOnError)

	branchName := walkFlag.String("b", "", "Branch to walk to")
	commitHash := walkFlag.String("c", "", "Commit hash to start walking from")

	return walkFlag, branchName, commitHash
}

func isValidBranchName(branchName *string) bool {
	return branchName != nil && len(*branchName) > 0
}

func isValidCommitHash(commitHash *string) bool {
	return commitHash != nil && len(*commitHash) == 40
}

func handleWalkCommand(walkFlag *flag.FlagSet, branchName *string, commitHash *string, projectPath string, args []string) {
	err := walkFlag.Parse(args)
	utils.Check(err, "Error parsing walk command arguments")

	validHash := isValidCommitHash(commitHash)
	validBranch := isValidBranchName(branchName)

	if validHash == validBranch { // XOR logic: one must be valid, but not both
		log.Fatal("walk command requires either a valid commit hash or a branch name but not both.\n")
	}

	walkCommand(projectPath, commitHash, branchName)

	if validBranch {
		log.Println("walking to ", *branchName)
	} else if validHash {
		log.Println("You are currently in detached HEAD state. Current HEAD position:", *commitHash)
	}
}

func CommandBuilder() commands.CommandBuilderOutput {
	command := "walk"
	walkFlag, branchName, commitHash := walkCliCommand(command)

	return commands.CommandBuilderOutput{
		Cmd: func(projectPath string, args []string) {
			handleWalkCommand(walkFlag, branchName, commitHash, projectPath, args)
		},
		Name: command,
	}
}
