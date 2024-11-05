package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/cmd/api"
	contentRepository "github.com/codevault-llc/humblebrag-api/internal/contents/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/core"
	"github.com/codevault-llc/humblebrag-api/internal/core/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	networkRepository "github.com/codevault-llc/humblebrag-api/internal/network/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")

	//user:password@(localhost:3306)/database_name
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5434/%s?sslmode=disable", postgresUser, postgresPassword, postgresDb)
	db, err := database.InitPostgres(connStr)
	if err != nil {
		log.Error("Error connecting to database", zap.Error(err))
		return
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

func SetupDatabases(db *sqlx.DB) {
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
