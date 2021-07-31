package test

import (
	"context"
	"testing"

	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/services/impl"
	"github.com/istt/api_gateway/pkg/fiber/instances"
	"github.com/istt/api_gateway/pkg/fiber/shared"
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
	collections, err := app.MongoDB.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)

	}
	t.Log(collections)

	userCollection := app.MongoDB.Collection("user")
	result := &shared.ManagedUserDTO{}
	var doc interface{}

	if err := bson.UnmarshalExtJSON([]byte(`{ "login": "admin" }`), true, &doc); err != nil {
		// handle error
		t.Fatal(err)
	}
	if err := userCollection.FindOne(context.TODO(), doc).Decode(result); err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", result.Login)
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
