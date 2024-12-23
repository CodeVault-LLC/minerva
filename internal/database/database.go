package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/codevault-llc/minerva/config"
	contentSchema "github.com/codevault-llc/minerva/internal/contents/models"
	coreSchema "github.com/codevault-llc/minerva/internal/core/models"
	networkSchema "github.com/codevault-llc/minerva/internal/network/models"
	"github.com/codevault-llc/minerva/pkg/logger"
)

type Database struct {
	Db clickhouse.Conn
}

// validateDatabase checks if the ClickHouse database is accessible.
func validateDatabase() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Config.DatabaseAddr},
		Auth: clickhouse.Auth{
			Username: config.Config.DatabaseUser,
			Password: config.Config.DatabasePass,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	if err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config.DatabaseName)); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

func NewDatabase() (*Database, error) {
	validateDatabase()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Config.DatabaseAddr},
		Auth: clickhouse.Auth{
			Database: config.Config.DatabaseName,
			Username: config.Config.DatabaseUser,
			Password: config.Config.DatabasePass,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "codevault-fingerprint", Version: "0.1"},
			},
		},
		// Disable TLS if not required
		TLS: nil,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
		Debugf:      log.Printf,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Ping the database to ensure the connection is established
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	db := &Database{Db: conn}

	if err := db.createSchema(); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	logger.Log.Info("Connected to ClickHouse")

	return db, nil
}

func (d *Database) createSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, schema := range coreSchema.CoreSchema {
		if err := d.Db.Exec(ctx, schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	for _, schema := range networkSchema.NetworkSchema {
		if err := d.Db.Exec(ctx, schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	for _, schema := range contentSchema.ContentSchema {
		if err := d.Db.Exec(ctx, schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	return nil
}

func (d *Database) GetDatabase() clickhouse.Conn {
	return d.Db
}

func (d *Database) Close() {
	d.Db.Close()
}
