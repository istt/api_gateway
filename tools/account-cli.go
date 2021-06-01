package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/internal/app/api-gateway/services/impl"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

var configFile, username, email, password, authorities, prefix string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "config", "configs/api-gateway.yaml", "API Gateway configuration file")
	flag.StringVar(&username, "username", "admin", "Username to create")
	flag.StringVar(&email, "email", "admin@localhost", "Email to create")
	flag.StringVar(&password, "password", "admin", "Password to create")
	flag.StringVar(&authorities, "authorities", "ROLE_USER,ROLE_ADMIN", "Authorities of user")
	flag.StringVar(&prefix, "prefix", "U:", "BuntDB user prefix to use")
	flag.Parse()

	// 1 - set default settings for components.
	app.MongoDBConfig()

	// 2 - override defaults with configuration file and watch changes
	app.ConfigInit(configFile)

	// 3 - bring up components
	app.MongoDBInit()

	instances.UserService = impl.NewUserServiceMongodb()

	var err error
	hash, err := instances.UserService.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	err = instances.UserService.SaveAccount(context.TODO(), &shared.ManagedUserDTO{
		UserDTO: shared.UserDTO{
			Id:          username,
			Login:       username,
			Email:       email,
			LangKey:     "en",
			Activated:   true,
			Authorities: []string{"ROLE_USER", "ROLE_ADMIN"},
		},
		CreatedBy:        "system",
		CreatedDate:      time.Now().Local().Format("2006-01-02"),
		Password:         hash,
		LastModifiedBy:   "system",
		LastModifiedDate: time.Now().Local().Format("2006-01-02"),
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully create user with ID: %s", username)
}
