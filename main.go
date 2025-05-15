package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"nit/commands"
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

func robe() {
	hash, header := calculateBlobHeader("what is up, doc?")
	gzipd := toGZip(header)
	fmt.Print(gzipd)

	if hash != "bd9dbf5aae1a3862dd1526723246b20206e5fc37\n" {
		fmt.Println("Error: ", hash)
	} else {
		fmt.Println("Success: ", hash)
	}
}

func main() {
	commands.InitCommand("nit")
}
