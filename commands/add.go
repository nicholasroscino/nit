package commands

import (
	"errors"
	"log"
	"nit/utils"
	"os"
	"strings"
)

type StagedObject struct {
	Hash      string
	Path      string
	Timestamp string
}

func getIndex(nitPath string) (map[string]string, error) {
	indexPath := nitPath + "/index"
	_, err := os.Stat(indexPath)

	if os.IsNotExist(err) {
		return nil, errors.New("index file does not exist")
	}

	content, err := os.ReadFile(indexPath)
	utils.Check(err, "Error reading the index file")

	index := make(map[string]string)

	lines := string(content)

	parts := strings.Split(lines, "\n")

	for _, line := range parts {
		if strings.Trim(line, " ") == "" {
			continue
		}

		arrayLine := strings.Split(line, " ")

		if len(arrayLine) != 2 {
			log.Fatal("Error parsing the index file")
		}

		index[arrayLine[0]] = strings.Trim(arrayLine[1], " ")
	}

	return index, nil
}

func writeIndex(nitPath string, index map[string]string) error {
	indexPath := nitPath + "/index"
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	utils.Check(err, "Error opening the index file")

	for key, value := range index {
		fileLineEntry := key + " " + strings.Trim(value, " ") + "\n"
		log.Println(fileLineEntry)
		_, err := file.WriteString(fileLineEntry)
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

func AddCommand(nitPath string, fileFullPath string) error {
	index, err := getIndex(nitPath)

	if err != nil {
		log.Println(err.Error())
		log.Println("recreating a new index file")
		index = make(map[string]string)
		err2 := os.WriteFile(""+nitPath+"/index", []byte{}, 0644)

		if err2 != nil {
			return errors.New("error creating the index file")
		}
	}

	hash, gzipContent := GetHashObject(fileFullPath)

	if val, ok := index[fileFullPath]; ok && val == hash {
		return errors.New("file already added to the index")
	}

	exist := utils.ObjectExist(nitPath, hash)
	if exist == false {
		SaveHashToFile(nitPath, hash, gzipContent)
	}

	index[fileFullPath] = hash

	err = writeIndex(nitPath, index)

	if err != nil {
		return err
	}

	return nil
}
