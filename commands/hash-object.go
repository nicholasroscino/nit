package commands

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"nit/utils"
	"os"
)

func toGZip(content string) string {
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

func calculateBlobHeader(myString string) (string, string) {
	header := fmt.Sprintf("blob %d\u0000", len(myString))
	h := sha1.New()

	fullBlob := header + myString
	done, err := io.WriteString(h, fullBlob)
	if err != nil {
		log.Fatal(err)
	}

	if done != len(header)+len(myString) {
		log.Fatal("Error writing to hash")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), fullBlob
}

func GetHashObject(fileFullPath string) (string, string) {
	fileContent, err := os.ReadFile(fileFullPath)
	utils.Check(err, "The file specified does not exist\n")

	hash, header := calculateBlobHeader(string(fileContent))
	gzipd := toGZip(header)

	return hash, gzipd
}

func SaveHashToFile(nitFolder string, hash string, gzipd string) {
	err := os.Mkdir(nitFolder+"/objects/"+hash[0:2], 0755)
	utils.Check(err, "Unable to create object directory with hash:"+hash[0:2]+"\n")

	err = os.WriteFile(nitFolder+"/objects/"+hash[0:2]+"/"+hash[2:], []byte(gzipd), 0644)
	utils.Check(err, "Unable to write object file\n")
}

func HashObjectCommand(nitFolder string, fileFullPath string) {
	//hash, gzipd := GetHashObject(fileFullPath)
	//SaveHashToFile(nitFolder, hash, gzipd)
}
