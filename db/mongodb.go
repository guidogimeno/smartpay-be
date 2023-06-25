package db

import (
	"context"
	"fmt"
	"log"

	"github.com/guidogimeno/smartpay-be/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	smartpayDb     = "smartpay"
	bcraCollection = "bcra"
)

type Mongo struct {
	client *mongo.Client
}

func NewMongo(port string) (*Mongo, error) {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + port)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return &Mongo{
		client: client,
	}, nil
}

func (m *Mongo) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected from MongoDB!")
}

func (m *Mongo) Create(f *types.FinancialData) error {
	collection := m.client.Database(smartpayDb).Collection(bcraCollection)
	fmt.Println(collection)

	what, err := collection.InsertOne(context.Background(), f)
	if err != nil {
		return err
	}

	fmt.Println("SE GUARDO", what)

	return nil
}

func (m *Mongo) Read() (*types.FinancialData, error) {
	collection := m.client.Database(smartpayDb).Collection(bcraCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*types.FinancialData
	for cursor.Next(context.Background()) {
		var data *types.FinancialData
		err := cursor.Decode(&data)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results[0], nil
}
