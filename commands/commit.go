package commands

import (
	"errors"
	"log"
	"nit/utils"
	"os"
	"strings"
)

type CommitObject struct {
	TreeHash string
	Parent   string
	Author   string
	Message  string
}

func serialiseCommitObject(commitObj CommitObject) string {
	return "tree " + commitObj.TreeHash + "\n" +
		"parent " + commitObj.Parent + "\n" +
		"author " + commitObj.Author + "\n" +
		"\n" +
		"message " + commitObj.Message + "\n"
}

func createCommitObject(nitPath string, treeHash string, author string, message string) (string, error) {
	// Take prev commit hash if exist at HEAD path
	headPath := nitPath + "/HEAD"

	currentHeadFilePathDesc, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	str := strings.Split(string(currentHeadFilePathDesc), " ")
	currentHeadFilePath := nitPath + "/" + str[1]

	fileContent, readHeadFileErr := os.ReadFile(currentHeadFilePath)

	if readHeadFileErr != nil {
		if !os.IsNotExist(readHeadFileErr) {
			return "", readHeadFileErr
		}
	}

	prevCommitHash := string(fileContent)

	newCommitObject := CommitObject{
		TreeHash: treeHash,
		Parent:   prevCommitHash,
		Author:   author,
		Message:  message,
	}
	commitFileContent := serialiseCommitObject(newCommitObject)
	hash, gzipd := utils.GetHashObjectFromContent(commitFileContent, "commit")
	utils.SaveHashToFileManaged(nitPath, hash, gzipd)

	// Write the new commit hash to HEAD
	err = os.WriteFile(string(currentHeadFilePath), []byte(hash), 0644)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CommitCommand(nitPath string, author string, message string) (string, error) {
	treeHash := WriteTreeCommand(nitPath)
	commitHash, err := createCommitObject(nitPath, treeHash, author, message)
	if err != nil {
		log.Println("Error creating commit object:", err)
		return "", errors.New("error creating commit object")
	}
	log.Println("Commit created with hash:", commitHash)
	return commitHash, nil
}
