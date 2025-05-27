package log

import (
	"nit/commands"
	"nit/commands/cat"
	"nit/utils"
)

func logCommand(projectPath string) error {
	// Get the nit path
	nitPath, err := utils.GetNitRepoFolder(projectPath)
	utils.Check(err, "This is not a nit repository")

	lastCommitHash, _, err := utils.GetLastCommitHash(nitPath)
	utils.Check(err, "No commits found in the repository.")

	for {
		content, err := cat.CatHeaderAndContent(nitPath, lastCommitHash)
		utils.Check(err, "Error reading the last commit object")

		commitObj := commands.NewCommitObject(content[1])
		println(commitObj.BeautyPrint())

		lastCommitHash = commitObj.Parent

		if lastCommitHash == "" {
			break
		}
	}

	return nil
}
