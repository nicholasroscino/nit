package init

import (
	. "nit/utils"
	"os"
)

func createNitFolders(path string) {
	var err = os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		Check(err, "Unable to create .nit directory\n")
	} else if os.IsExist(err) {
		Check(err, "Already in a nit repository\n")
	}

	err = os.Mkdir(path+"/objects", 0755)
	Check(err, "Unable to create .nit/objects directory\n")

	err = os.Mkdir(path+"/refs", 0755)
	Check(err, "Unable to create .nit/refs directory\n")

	err = os.Mkdir(path+"/refs/heads", 0755)
	Check(err, "Unable to create .nit/refs/heads directory\n")
}

func createNitFiles(path string) {
	var err = os.WriteFile(path+"/HEAD", []byte("ref: refs/heads/boss\n"), 0644)
	Check(err, "Unable to create .nit/HEAD file\n")

	d1 := []byte("[core]\n\trepositoryformatversion = 0\n")
	err = os.WriteFile(path+"/config", d1, 0644)
	Check(err, "Unable to create .nit/config file\n")

	err = os.WriteFile(path+"/description", []byte("Nit repository, write a description here\n"), 0644)
	Check(err, "Unable to create .nit/description file\n")
}

func initCommand(path string) {
	path = path + "/.nit"
	createNitFolders(path)
	createNitFiles(path)
}
