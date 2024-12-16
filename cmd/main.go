package main

import (
	"fmt"
	"log"

	orderModel "washit-api/app/order/model"
	userModel "washit-api/app/user/model"
	"washit-api/cmd/api"
	"washit-api/configs"
	dbs "washit-api/db"
	"washit-api/redis"
)

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		configs.Envs.DBUser,
		configs.Envs.DBPassword,
		configs.Envs.DBHost,
		configs.Envs.DBPort,
		configs.Envs.DBName,
	)

	db, err := dbs.NewDatabase(dsn)
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
