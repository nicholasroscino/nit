package commit

import (
	"errors"
	"log"
	"nit/commands"
	"nit/commands/cat"
	"nit/commands/write-tree"
	"nit/utils"
	"os"
	"strings"
)

func serialiseCommitObject(commitObj commands.CommitObject) string {
	output := "tree " + commitObj.TreeHash + "\n"

	if commitObj.Parent != "" {
		output += "parent " + commitObj.Parent + "\n"
	}

	output += "author " + commitObj.Author + "\n" +
		"\n" +
		"message " + commitObj.Message

	return output
}

func getTreeHashFromFileContent(fileContent []string) (string, error) {
	for _, line := range fileContent {
		if strings.HasPrefix(line, "tree ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "tree ")), nil
		}
	}
	return "", errors.New("no tree hash found in the commit object")
}

func createCommitObject(nitPath string, treeHash string, author string, message string) (string, error) {
	headPath := nitPath + "/HEAD"

	currentHeadFilePathDesc, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	str := strings.Split(string(currentHeadFilePathDesc), " ")
	currentHeadFilePath := nitPath + "/" + str[1]

	fileContent, readHeadFileErr := os.ReadFile(currentHeadFilePath)

	prevCommitTreeHash := ""

	if readHeadFileErr != nil && !os.IsNotExist(readHeadFileErr) {
		return "", readHeadFileErr
	}

	prevCommitHash := string(fileContent)

	if !os.IsNotExist(readHeadFileErr) {
		file, err2 := cat.CatHeaderAndContent(nitPath, prevCommitHash)

		if err2 != nil {
			log.Fatal("Error reading previous commit hash:", err2)
		}

		theType := strings.Split(file[0], " ")
		if len(theType) < 2 || theType[0] != "commit" {
			log.Fatal("The previous commit hash does not point to a valid commit object")
		}

		prevCommitTreeHash, err2 = getTreeHashFromFileContent(file)

		if err2 != nil {
			log.Fatal("Error getting tree hash from previous commit object:", err2)
		}
	}

	if prevCommitTreeHash != "" && prevCommitTreeHash == treeHash {
		log.Println("ma qui ci finisco?")
		log.Fatal("The tree hash is the same as the previous commit, no new commit created.")
	}

	newCommitObject := commands.CommitObject{
		TreeHash: treeHash,
		Parent:   prevCommitHash,
		Author:   author,
		Message:  message,
	}
	commitFileContent := serialiseCommitObject(newCommitObject)
	hash, gzipd := utils.GetHashObjectFromContent(commitFileContent, "commit")
	utils.SaveHashToFileManaged(nitPath, hash, gzipd)

	err = os.WriteFile(currentHeadFilePath, []byte(hash), 0644)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func commitCommand(projectPath string, author string, message string) (string, error) {
	nitPath, err := utils.GetNitRepoFolder(projectPath)
	utils.Check(err, "This is not a nit repository")

	treeHash := write_tree.WriteTreeCommand(projectPath)
	commitHash, err := createCommitObject(nitPath, treeHash, author, message)
	if err != nil {
		log.Println("Error creating commit object:", err)
		return "", errors.New("error creating commit object")
	}
	log.Println("Commit created with hash:", commitHash)
	return commitHash, nil
}
