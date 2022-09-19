package main

import (
	"kamari.com/db"
	"kamari.com/server"
	"log"
	"os"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Println("Shutting down server: ", err.Error())
		os.Exit(1)
	}

	server.Init()
}
