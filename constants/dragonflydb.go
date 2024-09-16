package constants

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitDragonflyDB() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test the connection
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to DragonflyDB: %v", err)
	}

	log.Println("Connected to DragonflyDB")
}
