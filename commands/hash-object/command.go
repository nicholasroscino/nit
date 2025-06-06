package hash_object

import (
	"nit/utils"
	"os"
)

func GetHashObject(fileFullPath string) (string, string) {
	fileContent, err := os.ReadFile(fileFullPath)
	utils.Check(err, "The file specified does not exist\n")
	hash, gzipd := utils.GetHashObjectFromContent(string(fileContent), "blob")

	return hash, gzipd
}

func hashObjectCommand(projectPath string, fileFullPath string) {
	nitFolder, err := utils.GetNitRepoFolder(projectPath)
	utils.Check(err, "This is not a nit repository")

	hash, gzipd := GetHashObject(fileFullPath)
	utils.SaveHashToFileManaged(nitFolder, hash, gzipd)
}
