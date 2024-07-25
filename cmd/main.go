package main

import (
	"log"

	"github.com/joho/godotenv"

	"go_bank/internal/route"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	route.InitializeRoute()
}
