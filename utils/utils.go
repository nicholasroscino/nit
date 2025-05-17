package utils

import (
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
