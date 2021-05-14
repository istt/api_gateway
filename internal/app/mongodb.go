package app

import (
	"context"
	"fmt"
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
		"mongodb.user":     "root",
		"mongodb.password": "",
		"mongodb.host":     "localhost27017",
		"mongodb.name":     "api-gateway",
	}, "."), nil)

}

// MongoDBInit create mongodb connection and assign it back to var MongoDB above
func MongoDBInit() {
	var err error

	mongodsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		Config.MustString("mongodb.user"),
		Config.String("mongodb.password"),
		Config.MustString("mongodb.host"),
	)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodsn))
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
