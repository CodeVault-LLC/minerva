package main

import (
	"log"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/cache"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
		log.Println(err)
	}

	parseUrl := os.Getenv("DATABASE_URL")
	stripeKey := os.Getenv("STRIPE_SECRET")
	stripe.Key = stripeKey

	postgres, err := database.InitPostgres(parseUrl)
	if err != nil {
		log.Println("Error connecting to database")
		log.Println(err)
	}

	redis, err := database.InitRedis()
	if err != nil {
		log.Println("Error connecting to redis")
		log.Println(err)
	}

	cache := cache.InitSessionManager()

	go updater.StartAutoUpdate(3 * time.Second)
	api.Start(postgres, redis, cache)
}
