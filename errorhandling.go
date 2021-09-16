package main

import (
	"log"
)

func logAndExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func logAndExitOnErrorWithAdditionalAction(err error, f func(error)) {
	if err != nil {
		f(err)
		log.Fatal(err)
	}
}

func logMessageAndExitOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
