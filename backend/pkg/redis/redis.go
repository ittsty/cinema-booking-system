package redis

import (
	"context"
	"log"

	"cinema-booking/pkg/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() {

	Client = redis.NewClient(&redis.Options{
		Addr: config.App.RedisAddr,
	})
	if err := Client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("failed to connect redis:", err)
	}

	log.Println("Redis connected")
}
