package commands

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

func fromGzip(content bytes.Buffer) (string, error) {
	zr, err := gzip.NewReader(&content)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer = bytes.Buffer{}
	if _, err := io.Copy(&out, zr); err != nil {
		log.Fatal(err)
	}

	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}

	return out.String(), nil
}

func CatFileCommand(projectFolder string, hash string) (string, error) {
	nitFolder := projectFolder + "/.nit"

	if len(hash) != 40 {
		return "", errors.New("the hash provided is not a valid SHA1 hash")
	}

	folder := nitFolder + "/objects/" + hash[0:2]
	file := folder + "/" + hash[2:]

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", errors.New("the hash provided does not exist in the repository")
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return "", errors.New("Error reading the file:" + file)
	}

	// Decompress the content
	var buf bytes.Buffer = bytes.Buffer{}
	buf.Write(content)

	theContent, err := fromGzip(buf)

	if err != nil {
		return "", errors.New("error decompressing the file provided: " + file)
	}

	return strings.Split(theContent, "\u0000")[1], nil
}
