package utils

import (
	"errors"
	"log"
	"os"
)

func Check(e error, msg string) {
	val := os.Getenv("NIT_DEBUG")

	if e != nil {
		if val == "1" {
			log.Fatal(msg, e.Error())
			return
		}

		log.Fatal(msg)
	}
}

func GetNitRepoFolder(path string) (string, error) {
	if _, err := os.Stat(path + "/.nit"); os.IsNotExist(err) {
		return "", errors.New(path + " is not a nit repository")
	}

	return path + "/.nit", nil
}

func ObjectExist(nitPath string, hash string) bool {
	if len(hash) != 40 {
		return false
	}

	pathToCheck := nitPath + "/objects/" + hash[0:2] + "/" + hash[2:]

	if _, err := os.Stat(pathToCheck); os.IsNotExist(err) {
		return false
	}

	return true
}
