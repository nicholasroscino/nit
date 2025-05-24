package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func GetNitFolder(projectDir string) string {
	return projectDir + "/.nit"
}

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

func ToGZip(content string) string {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(content))
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func CalculateHeader(myString string, fileType string) (string, string) {
	header := fmt.Sprintf(fileType+" %d\u0000", len(myString))
	h := sha1.New()

	fullFile := header + myString
	done, err := io.WriteString(h, fullFile)
	if err != nil {
		log.Fatal(err)
	}

	if done != len(header)+len(myString) {
		log.Fatal("Error writing to hash")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), fullFile
}

func SaveHashToFileManaged(nitFolder string, hash string, gzipd string) {
	err := os.Mkdir(nitFolder+"/objects/"+hash[0:2], 0755)
	Check(err, "Unable to create object directory with hash:"+hash[0:2]+"\n")

	err = os.WriteFile(nitFolder+"/objects/"+hash[0:2]+"/"+hash[2:], []byte(gzipd), 0644)
	Check(err, "Unable to write object file\n")
}

func SaveHashToFile(nitFolder string, hash string, gzipd string) error {
	err := os.Mkdir(nitFolder+"/objects/"+hash[0:2], 0755)
	if err != nil && os.IsExist(err) {
		return &HashAlreadyExist{}
	}

	err = os.WriteFile(nitFolder+"/objects/"+hash[0:2]+"/"+hash[2:], []byte(gzipd), 0644)
	if err != nil {
		return errors.New("Unable to write object file: " + err.Error())
	}

	return nil
}

func GetHashObjectFromContent(fileContent string, fileType string) (string, string) {
	hash, header := CalculateHeader(fileContent, fileType)
	gzipd := ToGZip(header)

	return hash, gzipd
}

//func GetNitRepoFolder(path string) (string, error) {
//	if _, err := os.Stat(path + "/.nit"); os.IsNotExist(err) {
//		return "", errors.New(path + " is not a nit repository")
//	}
//
//	return path + "/.nit", nil
//}

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
