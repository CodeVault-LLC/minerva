package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/minerva/cmd/api"
	"github.com/codevault-llc/minerva/config"
	"github.com/codevault-llc/minerva/internal/core"
	"github.com/codevault-llc/minerva/internal/core/models/repository"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/fingerprint"
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

	_, err = config.NewConfig()
	if err != nil {
		log.Error("Error loading config", zap.Error(err))
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

	setupScanning(db)

	if err := setupServices(); err != nil {
		log.Error("Error starting services", zap.Error(err))
		os.Exit(1)
	}

	go updater.StartAutoUpdate(20 * time.Minute)
	api.Start()
}

func setupServices() error {
	fingerprintClient, err := fingerprint.NewClient(config.Config.FingerprintServiceAddress)
	if err != nil {
		return err
	}

	fingerprint.FingerprintClient = fingerprintClient

	return nil
}

func setupScanning(db *database.Database) {
	pageAnalysis := core.NewPageAnalysis()

	repository.ScanRepository = repository.NewScanRepository(db)

	core.Scheduler = core.NewTaskScheduler(10)
	core.InspectorCore = core.NewInspector(db, pageAnalysis)
	core.Scheduler.Start(core.InspectorCore)
}
