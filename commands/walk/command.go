package walk

import (
	"errors"
	"log"
	"nit/commands"
	"nit/commands/cat"
	"nit/utils"
	"strings"
)

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

func walkAndWrite(currentTreeHash string, currentPath string, nitPath string, projectPath string) error {
	treeContent, err := cat.CatHeaderAndContent(nitPath, currentTreeHash)
	utils.Check(err, "Error reading the tree object")

	fileTypeAndSize := strings.Split(treeContent[0], " ")

	if fileTypeAndSize[0] != "tree" {
		log.Println("apparently not a tree :thinking: " + fileTypeAndSize[0])
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

func walkCommand(projectPath string, commitHash *string) error {
	nitPath, err := utils.GetNitRepoFolder(projectPath)
	utils.Check(err, "This is not a nit repository")

	headerAndContent, err := cat.CatHeaderAndContent(nitPath, *commitHash)
	utils.Check(err, "Error reading the commit object")

	headerWithSize := strings.Split(headerAndContent[0], " ")

	if headerWithSize[0] != "commit" {
		return errors.New("the provided hash does not point to a valid commit object")
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
			return err
		}

		fullPath := projectPath + "/" + stagedObject.Path

		err = utils.DeleteFile(fullPath)

		if err != nil {
			rollbackDeletedFiles(projectPath, deletedFiles)
			return err
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
		return err
	}

	return nil
}
