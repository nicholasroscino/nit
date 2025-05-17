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
	io.WriteString(h, header+myString)
	return fmt.Sprintf("%x\n", h.Sum(nil)), header
}

func HashObjectCommand(projectBasePath string, filePath string) {
	nitFolder := projectBasePath + "/.nit"
	fileFullPath := projectBasePath + "/" + filePath

	fileContent, err := os.ReadFile(fileFullPath)
	utils.Check(err, "The file specified does not exist\n")

	hash, header := calculateBlobHeader(string(fileContent))
	gzipd := toGZip(header)

	err = os.Mkdir(nitFolder+"/objects/"+hash[0:2], 0755)
	utils.Check(err, "Unable to create object directory with hash:"+hash[0:2]+"\n")

	err = os.WriteFile(nitFolder+"/objects/"+hash[0:2]+"/"+hash[2:], []byte(gzipd), 0644)
	utils.Check(err, "Unable to write object file\n")
}
