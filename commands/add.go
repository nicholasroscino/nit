package commands

import (
	"errors"
	"log"
	"nit/utils"
	"os"
	"strings"
	"time"
)

type StagedObject struct {
	Hash      string
	Path      string
	Timestamp string
}

func serializeStagedObject(stagedObject StagedObject) string {
	return stagedObject.Hash + " " + stagedObject.Path + " " + stagedObject.Timestamp
}

func parseStagedObject(line string) (StagedObject, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return StagedObject{}, errors.New("error parsing the staged object")
	}

	return StagedObject{
		Hash:      strings.Trim(parts[0], " "),
		Path:      strings.Trim(parts[1], " "),
		Timestamp: strings.Trim(parts[2], " "),
	}, nil
}

func getIndex(nitPath string) (map[string]StagedObject, error) {
	indexPath := nitPath + "/index"
	_, err := os.Stat(indexPath)

	if os.IsNotExist(err) {
		return nil, errors.New("index file does not exist")
	}

	content, err := os.ReadFile(indexPath)
	utils.Check(err, "Error reading the index file")

	index := make(map[string]StagedObject)

	lines := string(content)

	parts := strings.Split(lines, "\n")

	for _, line := range parts {
		if strings.Trim(line, " ") == "" {
			continue
		}

		stagedObj, err := parseStagedObject(line)

		if err != nil {
			return nil, err
		}

		index[stagedObj.Path] = StagedObject{
			Hash:      stagedObj.Hash,
			Path:      stagedObj.Path,
			Timestamp: stagedObj.Timestamp,
		}
	}

	return index, nil
}

func writeIndex(nitPath string, index map[string]StagedObject) error {
	indexPath := nitPath + "/index"
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	utils.Check(err, "Error opening the index file")

	for _, value := range index {
		fileLineEntry := serializeStagedObject(value)
		_, err := file.WriteString(fileLineEntry + "\n")
		if err != nil {
			return errors.New("error writing to the index file")
		}
	}
	err = file.Close()

	if err != nil {
		return errors.New("error closing the index file")
	}

	return nil
}

func AddCommand(nitPath string, filePath string) error {
	index, err := getIndex(nitPath)

	projectPath := nitPath[0:strings.LastIndex(nitPath, ".nit")]

	fileFullPath := projectPath + filePath

	if err != nil {
		log.Println(err.Error())
		log.Println("recreating a new index file")
		index = make(map[string]StagedObject)
		err2 := os.WriteFile(""+nitPath+"/index", []byte{}, 0644)

		if err2 != nil {
			return errors.New("error creating the index file")
		}
	}

	hash, gzipContent := GetHashObject(fileFullPath)

	if val, ok := index[filePath]; ok && val.Hash == hash {
		return errors.New("file already added to the index")
	}

	exist := utils.ObjectExist(nitPath, hash)
	if exist == false {
		SaveHashToFile(nitPath, hash, gzipContent)
	}

	index[filePath] = StagedObject{
		Hash:      hash,
		Path:      filePath,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err = writeIndex(nitPath, index)

	if err != nil {
		return err
	}

	return nil
}
