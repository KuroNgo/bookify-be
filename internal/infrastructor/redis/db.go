package redis

import (
	"bookify/internal/config"
	"bookify/pkg/shared/helper"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func NewRedisClient(env *config.Database) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var redisAddr string
	if env.DBRUser != "" && env.DBRPassword != "" {
		redisAddr = fmt.Sprintf("redis://%s:%s@%s:%s", env.DBRUser, env.DBRPassword, env.DBRHost, env.DBRPort)
	} else {
		redisAddr = fmt.Sprintf("%s:%s", env.DBRHost, env.DBRPort)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: env.DBRPassword,
		DB:       0, // Sử dụng database 0 mặc định
	})

	if err := client.Ping(ctx).Err(); err != nil {
		helper.FailToError(err, "error while trying to ping redis:")
		return nil
	}

	log.Println("Connected to Redis!")
	return client
}

func CloseRedisConnection(client *redis.Client) {
	if err := client.Close(); err != nil {
		log.Fatal("error while closing Redis connection:", err)
	}

	log.Println("Connection to Redis closed.")
}
