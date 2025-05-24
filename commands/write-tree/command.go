package write_tree

import (
	"log"
	"nit/commands"
	"nit/utils"
	"strings"
)

func popFirstPathReturnRemaining(path string) (string, string) {
	// Pop the first path and return the remaining path
	if strings.Contains(path, "/") {
		firstPath := path[0:strings.Index(path, "/")]
		remainingPath := path[strings.Index(path, "/")+1:]
		return firstPath, remainingPath
	}

	return path, ""
}

func printTree(node *commands.NitNode, indent string) {
	if node == nil {
		return
	}

	log.Printf("%s%s %s(%s)\n", indent, node.Id, node.Hash, node.Type)
	for _, file := range node.Files {
		printTree(file, indent+"  ")
	}
}

func CreateTreeFile(nitPath string, node *commands.NitNode) (string, []commands.FileContent) {
	arr := make([]commands.FileContent, 0)

	if node.Type == "blob" {
		arr := make([]commands.FileContent, 0)
		arr = append(arr, commands.FileContent{
			Id:   node.Id,
			Hash: node.Hash,
			Type: "blob",
		})
		return node.Hash, arr
	}

	filesOfPath := make([]commands.FileContent, 0)
	for _, file := range node.Files {
		_, files := CreateTreeFile(nitPath, file)
		for _, f := range files {
			filesOfPath = append(filesOfPath, f)
		}
	}

	fileContent := ""
	for _, file := range filesOfPath {
		fileContent += file.Type + " " + file.Hash + " " + file.Id + "\n"
	}

	hash, gzipd := utils.GetHashObjectFromContent(fileContent, "tree")
	saveFileError := utils.SaveHashToFile(nitPath, hash, gzipd)

	if saveFileError != nil && utils.IsHashAlreadyExist(saveFileError) {
		log.Fatal("nothing to commit, working tree clean")
	}

	arr = append(arr, commands.FileContent{
		Id:   node.Id,
		Hash: hash,
		Type: "tree",
	})

	return hash, arr
}

func WriteTree(currentPath *utils.StagedObject, rootNode *commands.NitNode) *commands.NitNode {
	if currentPath.Path == "" {
		return rootNode
	}

	head, tail := popFirstPathReturnRemaining(currentPath.Path)

	var child *commands.NitNode
	for i := range rootNode.Files {
		if rootNode.Files[i].Id == head {
			child = rootNode.Files[i]
			break
		}
	}

	if child == nil {
		newType := "tree"
		if tail == "" {
			newType = "blob"
		}

		newNode := commands.NitNode{
			Id:    head,
			Type:  newType,
			Files: []*commands.NitNode{},
			Hash:  "",
		}

		if newType == "blob" {
			newNode.Hash = currentPath.Hash
		}

		rootNode.Files = append(rootNode.Files, &newNode)
		child = &newNode
	}

	if tail != "" {
		WriteTree(&utils.StagedObject{
			Hash:      currentPath.Hash,
			Path:      tail,
			Timestamp: currentPath.Timestamp,
		}, child)
	}

	return rootNode
}

func WriteTreeCommand(projectPath string) string {
	nitPath := utils.GetNitFolder(projectPath)

	index, err := utils.GetIndex(nitPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	rootNode := &commands.NitNode{
		Id:    projectPath,
		Hash:  "",
		Files: make([]*commands.NitNode, 0),
		Type:  "tree",
	}

	for _, val := range index {
		WriteTree(&val, rootNode)
	}

	hash, _ := CreateTreeFile(nitPath, rootNode)

	return hash
}
