package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mainawycliffe/kamanda/cmd"
)

func main() {
	// this is temporary way of passing the credentials for a firebase project
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
