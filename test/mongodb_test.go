package test

import (
	"context"
	"testing"

	"github.com/istt/api_gateway/internal/app"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDB(t *testing.T) {

	app.MongoDBConfig()

	app.MongoDBInit()

	// try access mongodb to do something
	dbnames, err := app.MongoClient.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dbnames)

	// test db
	app.MongoDB.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)

	}
	t.Log(dbnames)
}
