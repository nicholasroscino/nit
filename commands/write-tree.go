package commands

import (
	"log"
	"nit/utils"
	"strings"
)

type NitNode struct {
	Hash  string
	Files []NitNode
	Type  string
}

func isFolder(path string) bool {
	// Check if the path is a folder
	if strings.Contains(path, "/") {
		return true
	}

	return false
}

func popFirstPathReturnRemaining(path string) (string, string) {
	// Pop the first path and return the remaining path
	if strings.Contains(path, "/") {
		firstPath := path[0:strings.Index(path, "/")]
		remainingPath := path[strings.Index(path, "/")+1:]
		return firstPath, remainingPath
	}

	return path, ""
}

func WriteTree(currentPath *utils.StagedObject, rootNode *NitNode) []NitNode {
	arr := make([]NitNode, 0)
	if isFolder(currentPath.Path) {
		popped, remaining := popFirstPathReturnRemaining(currentPath.Path)

		arr = append(arr, NitNode{
			Hash: popped,
			Files: WriteTree(&utils.StagedObject{
				Hash:      currentPath.Hash,
				Path:      remaining,
				Timestamp: currentPath.Timestamp,
			}, rootNode),
			Type: "tree",
		})
		return arr
	}

	arr = append(arr, NitNode{
		Hash:  currentPath.Path,
		Files: []NitNode{},
		Type:  "blob",
	})
	return arr
}

func WriteTreeCommand(nitPath string) {
	index, err := utils.GetIndex(nitPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	folderMap := make(map[string]NitNode)

	for _, val := range index {
		println(folderMap, val)
	}
}
