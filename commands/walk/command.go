package walk

import (
	"errors"
	"log"
	"nit/commands"
	"nit/commands/cat"
	"nit/utils"
	"os"
	"strings"
)

func walkDetachedState(nitPath string, projectPath string, commitHash *string) {
	moveToCommit(nitPath, projectPath, *commitHash)
	utils.WriteOnHead(nitPath, *commitHash)
}

func moveToCommit(nitPath string, projectPath string, commitHash string) {
	headerAndContent, err := cat.CatHeaderAndContent(nitPath, commitHash)
	utils.Check(err, "Error reading the commit object")

	headerWithSize := strings.Split(headerAndContent[0], " ")

	if headerWithSize[0] != "commit" {
		log.Fatal("the provided hash does not point to a valid commit object")
	}

	lines, err := utils.GetIndexFile(nitPath)
	utils.Check(err, "Error reading the index file")

	deletedFiles := make([]struct {
		path string
		hash string
	}, 0)

	for _, line := range lines {
		if line == "" {
			continue
		}

		stagedObject, err := utils.ParseStagedObject(line)

		if err != nil {
			rollbackDeletedFiles(projectPath, deletedFiles)
			log.Fatal(err.Error())
		}

		fullPath := projectPath + "/" + stagedObject.Path

		err = utils.DeleteFile(fullPath)

		if err != nil {
			rollbackDeletedFiles(projectPath, deletedFiles)
			log.Fatal(err.Error())
		}

		deletedFiles = append(deletedFiles, struct {
			path string
			hash string
		}{path: stagedObject.Path, hash: stagedObject.Hash})
	}

	commitObj := commands.NewCommitObject(headerAndContent[1])

	err = walkAndWrite(commitObj.TreeHash, projectPath, nitPath, projectPath)

	if err != nil {
		rollbackDeletedFiles(projectPath, deletedFiles)
		log.Fatal(err.Error())
	}
}

func walkToBranch(nitPath string, projectPath string, branchName string) {
	branchFilePath := nitPath + "/refs/heads/" + branchName

	lastCommit, _, err := utils.GetLastCommitHash(nitPath)
	utils.Check(err, "Error getting the last commit hash")

	if _, err = os.Stat(branchFilePath); err != nil {
		log.Println("Branch does not exist:", branchName)

		err = os.WriteFile(branchFilePath, []byte(lastCommit), 0644)
		utils.Check(err, "Error creating the branch file")
		log.Println("Created a new branch with the name:", branchName)
	} else {
		branchContent, err := os.ReadFile(branchFilePath)
		utils.Check(err, "Error reading the branch file")

		fileContent := string(branchContent)

		if len(fileContent) != 40 {
			log.Fatal("Error reading the branch file: it does not contain a valid commit hash")
		}

		moveToCommit(nitPath, projectPath, fileContent)
	}

	utils.WriteOnHead(nitPath, "ref: "+"refs/heads/"+branchName)
}

func walkCommand(projectPath string, commitHash *string, branchName *string) {
	nitPath, err := utils.GetNitRepoFolder(projectPath)
	utils.Check(err, "This is not a nit repository")

	if commitHash != nil && *commitHash != "" {
		walkDetachedState(nitPath, projectPath, commitHash)
	} else if branchName != nil && *branchName != "" {
		walkToBranch(nitPath, projectPath, *branchName)
	}
}

func walkAndWrite(currentTreeHash string, currentPath string, nitPath string, projectPath string) error {
	treeContent, err := cat.CatHeaderAndContent(nitPath, currentTreeHash)
	utils.Check(err, "Error reading the tree object")

	fileTypeAndSize := strings.Split(treeContent[0], " ")

	if fileTypeAndSize[0] != "tree" {
		return errors.New("the provided hash does not point to a valid tree object")
	}

	treeContentStruct, err := getTreeNodeFromContent(treeContent[1])
	utils.Check(err, "Error parsing the tree content")

	for _, file := range treeContentStruct {
		if file.Type == "blob" {
			fullPath := currentPath + "/" + file.PathName

			err := cat.CatHashToFile(projectPath, file.Hash, fullPath)
			utils.Check(err, "Error restoring file: "+fullPath)

		} else if file.Type == "tree" {
			err := walkAndWrite(file.Hash, currentPath+"/"+file.PathName, nitPath, projectPath)
			utils.Check(err, "Error walking through tree object: "+file.PathName)
		} else {
			log.Println("Unknown type in tree object: " + file.Type)
		}
	}

	return nil
}

func getTreeNodeFromContent(treeContentRaw string) ([]*commands.NitTreeFileContent, error) {
	treeContent := strings.Split(treeContentRaw, "\n")

	treeContentStruct := make([]*commands.NitTreeFileContent, 0)

	for _, line := range treeContent {
		if strings.Trim(line, " ") == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			return nil, errors.New("invalid tree content format")
		}
		treeFile := &commands.NitTreeFileContent{
			Type:     parts[0],
			Hash:     parts[1],
			PathName: parts[2],
		}

		treeContentStruct = append(treeContentStruct, treeFile)
	}

	return treeContentStruct, nil

}

func rollbackDeletedFiles(projectPath string, deletedFiles []struct {
	path string
	hash string
}) {
	for _, file := range deletedFiles {
		fullPath := projectPath + "/" + file.path

		err := cat.CatHashToFile(projectPath, file.hash, fullPath)
		utils.Check(err, "Error restoring file: "+file.path)

	}
}
