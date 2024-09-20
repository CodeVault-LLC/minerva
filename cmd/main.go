package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/cache"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	log, err := logger.NewLogger(true, true, "humblebrag-api.log", logger.DEBUG)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	defer log.Close()

	err = godotenv.Load()

	if err != nil {
		log.Error("Error loading .env file", err)
	}

	parseUrl := os.Getenv("DATABASE_URL")
	stripeKey := os.Getenv("STRIPE_SECRET")
	stripe.Key = stripeKey
	config.InitDiscordAuth()

	postgres, err := database.InitPostgres(parseUrl)
	if err != nil {
		log.Error("Error connecting to database", err)
	}
	log.Info("Postgres connected")

	redis, err := database.InitRedis()
	if err != nil {
		log.Error("Error connecting to redis", err)
	}
	log.Info("Redis connected")

	cache := cache.InitSessionManager()

	go updater.StartAutoUpdate(3 * time.Second)
	api.Start(postgres, redis, cache)
}
