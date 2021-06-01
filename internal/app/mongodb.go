package app

import (
	"context"
	"log"
	"time"

	"github.com/knadh/koanf/providers/confmap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Create global mongoDB connection, so that other components can access with app.MongoDB.
	// var MongoDB
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

// MongoDBConfig provide default settings for MongoDB
func MongoDBConfig() {
	Config.Load(confmap.Provider(map[string]interface{}{
		"mongodb.username": "root",
		"mongodb.password": "",
		"mongodb.host":     "127.0.0.1",
		"mongodb.port":     27017,
		"mongodb.authdb":   "",
		"mongodb.name":     "test",
	}, "."), nil)

}

const (
	SchemeMongoDB    = "mongodb"
	SchemeMongoDBSRV = "mongodb+srv"
)

// MongoDBInit create mongodb connection and assign it back to var MongoDB above
func MongoDBInit() {
	var err error

	// mongodb := fmt.Sprintf("%s:%s@%s:%d/test?retryWrites=true&w=majority",
	// 	Config.MustString("mongodb.username"),
	// 	Config.String("mongodb.password"),
	// 	Config.MustString("mongodb.host"),
	// 	Config.MustInt("mongodb.port"),
	//
	// FIXME: using fmt.Sprintf
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	MongoClient = client

	MongoDB = client.Database(Config.MustString("mongodb.name"))
}
