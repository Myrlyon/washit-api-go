package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"

	"washit-api/app/user/model"
	"washit-api/cmd/api"
	"washit-api/configs"
	dbs "washit-api/db"
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

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate models", err)
	}

	initRedis(context.Background())

	server := api.NewServer(db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initRedis(ctx context.Context) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully.")

	return rdb, nil
}
