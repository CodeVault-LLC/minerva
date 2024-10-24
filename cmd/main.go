package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/cache"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
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
		log.Error("Error loading .env file %v", err)
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")

	stripeKey := os.Getenv("STRIPE_SECRET")
	stripe.Key = stripeKey
	config.InitDiscordAuth()

	postgres, err := database.InitPostgres(fmt.Sprintf("postgres://%s:%s@localhost:5434/%s?sslmode=disable", postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Error("Error connecting to database %v", err)
	}
	log.Info("Postgres connected")

	redis, err := database.InitRedis()
	if err != nil {
		log.Error("Error connecting to redis %v", err)
	}
	log.Info("Redis connected")

	cache := cache.InitSessionManager()

	websites.InitBrowser()

	go updater.StartAutoUpdate(20 * time.Minute)
	api.Start(postgres, redis, cache)
}
