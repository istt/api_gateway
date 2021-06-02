package test

import (
	"context"
	"testing"

	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/internal/app/api-gateway/services/impl"
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

func TestUser(t *testing.T) {
	app.MongoDBConfig()
	app.MongoDBInit()
	instances.UserService = impl.NewUserServiceMongodb()
	res, err := instances.UserService.GetUserByUsername(context.TODO(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}
