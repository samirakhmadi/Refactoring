package main

import (
	"log"
	"refactoring/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("error init app:", err)
	}
}
