package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/minerva/cmd/api"
	"github.com/codevault-llc/minerva/internal/core"
	"github.com/codevault-llc/minerva/internal/core/models/repository"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/updater"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/joho/godotenv"
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

	db, err := database.NewDatabase()
	if err != nil {
		log.Error("Error connecting to database %v", zap.Error(err))
	}

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

	SetupScanning(db)

	go updater.StartAutoUpdate(20 * time.Minute)
	api.Start()
}

func SetupScanning(db *database.Database) {
	core.InitializeBrowser()

	repository.ScanRepository = repository.NewScanRepository(db)

	core.Scheduler = core.NewTaskScheduler(10)
	core.InspectorCore = core.NewInspector(db)
	core.Scheduler.Start(core.InspectorCore)
}
