package utils

import (
	"errors"
	"os"
	"strings"
)

type StagedObject struct {
	Hash      string
	Path      string
	Timestamp string
}

func SerializeStagedObject(stagedObject StagedObject) string {
	return stagedObject.Hash + " " + stagedObject.Path + " " + stagedObject.Timestamp
}

func ParseStagedObject(line string) (StagedObject, error) {
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

func GetIndex(nitPath string) (map[string]StagedObject, error) {
	indexPath := nitPath + "/index"
	_, err := os.Stat(indexPath)

	if os.IsNotExist(err) {
		return nil, errors.New("index file does not exist")
	}

	content, err := os.ReadFile(indexPath)
	Check(err, "Error reading the index file")

	index := make(map[string]StagedObject)

	lines := string(content)

	parts := strings.Split(lines, "\n")

	for _, line := range parts {
		if strings.Trim(line, " ") == "" {
			continue
		}

		stagedObj, err := ParseStagedObject(line)

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

func WriteIndex(nitPath string, index map[string]StagedObject) error {
	indexPath := nitPath + "/index"
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	Check(err, "Error opening the index file")

	for _, value := range index {
		fileLineEntry := SerializeStagedObject(value)
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
