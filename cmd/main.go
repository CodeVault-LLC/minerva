package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	"github.com/codevault-llc/humblebrag-api/internal/core"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Initialize the logger
	log, err := logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}

	err = godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file", zap.Error(err))
	}

	// Initialize the database
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")

	db, err := database.InitPostgres(fmt.Sprintf("postgres://%s:%s@localhost:5434/%s?sslmode=disable", postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Error("Error connecting to database %v", zap.Error(err))
	}
	log.Info("Connected to database")

	// Initialize Redis
	_, err = database.InitRedis()
	if err != nil {
		log.Error("Error connecting to redis %v", zap.Error(err))
	}
	log.Info("Connected to redis")

	// Initialize AWS
	err = database.InitAWS()
	if err != nil {
		log.Error("Error connecting to AWS %v", zap.Error(err))
	}
	log.Info("Connected to AWS")

	SetupDatabases(db)
	SetupScanning()

	go updater.StartAutoUpdate(20 * time.Minute)
	api.Start()
}

func SetupDatabases(db *gorm.DB) {
	repository.ScanRepository = repository.NewScanRepository(db)
	repository.NetworkRepository = repository.NewNetworkRepository(db)
	repository.LicenseRepository = repository.NewLicenseRepository(db)
	repository.ContentRepository = repository.NewContentRepo(db)
	repository.FindingRepository = repository.NewFindingRepo(db)
	repository.DnsRepository = repository.NewDnsRepository(db)
	repository.WhoisRepository = repository.NewWhoisRepository(db)
	repository.CertificateRepository = repository.NewCertificateRepository(db)
	repository.MetadataRepository = repository.NewMetadataRepository(db)
	repository.ScreenshotRepository = repository.NewScreenshotRepo(db)
	repository.RedirectRepository = repository.NewRedirectRepo(db)
}

func SetupScanning() {
	core.InitializeBrowser()

	core.Scheduler = core.NewTaskScheduler(10)
	core.InspectorCore = core.NewInspector()
	core.Scheduler.Start(core.InspectorCore)
}
