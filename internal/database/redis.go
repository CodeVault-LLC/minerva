package database

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() (*redis.Client, error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:               "localhost:6379",
		IdleTimeout:        5 * time.Minute,
		DB:                 0,
		MaxRetries:         2,
		PoolSize:           100, // Increase the pool size
		MinIdleConns:       2,   // Maintain some minimum idle connections
		ReadTimeout:        100 * time.Second,
		WriteTimeout:       100 * time.Second,
		IdleCheckFrequency: 5 * time.Minute,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Error connecting to Redis")
		log.Println(err)

		return nil, err
	} else {
		log.Println("Connected to Redis")
	}

	return Rdb, nil
}
