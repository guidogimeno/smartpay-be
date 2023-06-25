package main

import (
	"log"
	"os"

	"github.com/guidogimeno/smartpay-be/api"
	"github.com/guidogimeno/smartpay-be/db"
	"github.com/guidogimeno/smartpay-be/worker"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbPort := os.Getenv("DB_PORT")

	mongo, err := db.NewMongo(dbPort)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Close()

	worker := worker.New(mongo)
	go worker.Start()

	apiServer := api.NewAPIServer(port, mongo)
	log.Fatal(apiServer.Run())
}
