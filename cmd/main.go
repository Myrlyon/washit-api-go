package main

import (
	"context"
	"fmt"
	"log"
	"washit-api/cmd/api"
	"washit-api/configs"
	dbs "washit-api/db"
	"washit-api/redis"
	"washit-api/utils"

	firebase "firebase.google.com/go"
	"github.com/go-playground/validator"
	"google.golang.org/api/option"
)

//	@title			Washit API
//	@version		1.0
//	@description	Swagger for washit app.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Marlen E. Satriani
//	@contact.email	marlendotedots@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	db, err := dbs.NewDatabase(configs.Envs.URI)
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	err = db.AutoMigrate(utils.ModelList...)
	if err != nil {
		log.Fatal("Failed to migrate models", err)
	}

	cache := redis.New(redis.Config{
		Address:  fmt.Sprintf("%s:%s", configs.Envs.RedisHost, configs.Envs.RedisPort),
		Password: configs.Envs.RedisPassword,
		Database: configs.Envs.RedisDB,
	})

	opt := option.WithCredentialsFile("washit-445307-firebase-adminsdk-ypdz8-c61d567af6.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	validate := validator.New()

	server := api.NewServer(validate, db, cache, app)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
