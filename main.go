package main

import (
	"log"
	"os"

	"github.com/guidogimeno/smartpay-be/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	apiServer := api.NewAPIServer(port)
	log.Fatal(apiServer.Run())
}
