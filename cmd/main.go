package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/minerva/cmd/api"
	contentRepository "github.com/codevault-llc/minerva/internal/contents/models/repository"
	"github.com/codevault-llc/minerva/internal/core"
	"github.com/codevault-llc/minerva/internal/core/models/repository"
	"github.com/codevault-llc/minerva/internal/database"
	networkRepository "github.com/codevault-llc/minerva/internal/network/models/repository"
	"github.com/codevault-llc/minerva/internal/updater"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	db, err := database.NewDatabase()
	if err != nil {
		log.Error("Error connecting to database %v", zap.Error(err))
	}

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

func SetupDatabases(db *database.Database) {
	repository.ScanRepository = repository.NewScanRepository(db)
	networkRepository.NetworkRepository = networkRepository.NewNetworkRepository(db)
	contentRepository.ContentRepository = contentRepository.NewContentRepo(db)
	contentRepository.FindingRepository = contentRepository.NewFindingRepo(db)
	networkRepository.DnsRepository = networkRepository.NewDnsRepository(db)
	networkRepository.WhoisRepository = networkRepository.NewWhoisRepository(db)
	networkRepository.CertificateRepository = networkRepository.NewCertificateRepository(db)
}

func SetupScanning() {
	core.InitializeBrowser()

	core.Scheduler = core.NewTaskScheduler(10)
	core.InspectorCore = core.NewInspector()
	core.Scheduler.Start(core.InspectorCore)
}
