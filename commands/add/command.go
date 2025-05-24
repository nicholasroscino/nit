package add

import (
	"errors"
	"log"
	"nit/commands/hash-object"
	"nit/utils"
	"os"
	"time"
)

func addCommand(projectPath string, filePath string) error {
	nitPath := utils.GetNitFolder(projectPath)

	index, err := utils.GetIndex(nitPath)
	fileFullPath := projectPath + "/" + filePath

	if err != nil {
		log.Println(err.Error())
		log.Println("recreating a new index file")
		index = make(map[string]utils.StagedObject)
		err2 := os.WriteFile(""+nitPath+"/index", []byte{}, 0644)

		if err2 != nil {
			return errors.New("error creating the index file")
		}
	}

	hash, gzipContent := hash_object.GetHashObject(fileFullPath)

	if val, ok := index[filePath]; ok && val.Hash == hash {
		return errors.New("file already added to the index")
	}

	exist := utils.ObjectExist(nitPath, hash)
	if exist == false {
		utils.SaveHashToFileManaged(nitPath, hash, gzipContent)
	}

	index[filePath] = utils.StagedObject{
		Hash:      hash,
		Path:      filePath,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err = utils.WriteIndex(nitPath, index)

	if err != nil {
		return err
	}

	return nil
}
