package main

import (
	"fmt"
	"log"

	"github.com/guidogimeno/smartpay-be/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	inflation, err := services.ScrapInflation("2023-01-01", "2023-05-01")
	if err != nil {
		fmt.Println(err)
	}

	for _, b := range inflation {
		fmt.Println(b)
	}

	// port := os.Getenv("PORT")

	// apiServer := api.NewAPIServer(port)
	// log.Fatal(apiServer.Run())
}
