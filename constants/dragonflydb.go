package constants

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitDragonflyDB() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Use the appropriate address and port
	})

	// Retry connection
	for i := 0; i < 5; i++ {
		_, err := Rdb.Ping(Ctx).Result()
		if err == nil {
			log.Println("Connected to DragonflyDB")
			return
		}
		log.Printf("Failed to connect to DragonflyDB (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("Could not connect to DragonflyDB after 5 attempts")
}
