package cat

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

func CatHeaderAndContent(nitFolder string, hash string) ([]string, error) {
	if len(hash) != 40 {
		return nil, errors.New("the hash provided is not a valid SHA1 hash")
	}

	folder := nitFolder + "/objects/" + hash[0:2]
	file := folder + "/" + hash[2:]

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.New("the hash provided does not exist in the repository")
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New("Error reading the file:" + file)
	}

	var buf = bytes.Buffer{}
	buf.Write(content)

	theContent, err := fromGzip(buf)

	if err != nil {
		return nil, errors.New("error decompressing the file provided: " + file)
	}

	return strings.Split(theContent, "\u0000"), nil
}

func catFileCommand(projectFolder string, hash string) (string, error) {
	nitFolder := projectFolder + "/.nit"
	file, err := CatHeaderAndContent(nitFolder, hash)

	if err != nil {
		return "", err
	}

	if len(file) < 2 {
		return "", errors.New("the file provided is not a valid nit file")
	}

	return file[1], nil
}
