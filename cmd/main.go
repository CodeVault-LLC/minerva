package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	"github.com/codevault-llc/humblebrag-api/internal/core"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}

	err = godotenv.Load()

	if err != nil {
		log.Error("Error loading .env file", zap.Error(err))
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")

	_, err = database.InitPostgres(fmt.Sprintf("postgres://%s:%s@localhost:5434/%s?sslmode=disable", postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Error("Error connecting to database %v", zap.Error(err))
	}
	log.Info("Connected to database")

	_, err = database.InitRedis()
	if err != nil {
		log.Error("Error connecting to redis %v", zap.Error(err))
	}
	log.Info("Connected to redis")

	err = database.InitAWS()
	if err != nil {
		log.Error("Error connecting to AWS %v", zap.Error(err))
	}
	log.Info("Connected to AWS")

	websites.InitializeBrowser()

	core.Scheduler = core.NewTaskScheduler(10)

	go updater.StartAutoUpdate(20 * time.Minute)
	api.Start()
}
