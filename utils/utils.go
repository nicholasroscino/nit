package utils

import "log"

func Check(e error, msg string) {
	if e != nil {
		log.Fatal(msg, e.Error())
	}
}
