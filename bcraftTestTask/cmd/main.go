package main

import (
	"bcraftTestTask/internal/app"
	"log"
)

func main() {
	newApp, err := app.NewApp()

	if err != nil {
		log.Println(err)
		return
	}

	err = newApp.Start()

	if err != nil {
		log.Println(err)
		return
	}
}
