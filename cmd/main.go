package main

import (
	"fmt"
	"log"

	orderModel "washit-api/app/order/dto/model"
	userModel "washit-api/app/user/dto/model"
	"washit-api/cmd/api"
	"washit-api/configs"
	dbs "washit-api/db"
	"washit-api/redis"
)

func main() {
	db, err := dbs.NewDatabase(configs.Envs.URI)
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	err = db.AutoMigrate(&userModel.User{}, &orderModel.Order{})
	if err != nil {
		log.Fatal("Failed to migrate models", err)
	}

	cache := redis.New(redis.Config{
		Address:  fmt.Sprintf("%s:%s", configs.Envs.RedisHost, configs.Envs.RedisPort),
		Password: configs.Envs.RedisPassword,
		Database: configs.Envs.RedisDB,
	})

	server := api.NewServer(db, cache)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
