package main

import (
	"fmt"
	"log"

	"github.com/forthedreamers-server/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = database.ConnectDB()

	fmt.Println("for the dreamers")
}
