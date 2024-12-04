package database

import (
	"fmt"

	contentSchema "github.com/codevault-llc/minerva/internal/contents/models"
	coreSchema "github.com/codevault-llc/minerva/internal/core/models"
	networkSchema "github.com/codevault-llc/minerva/internal/network/models"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type Database struct {
	Db *gocql.Session
}

func NewDatabase() (*Database, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum

	tempSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer tempSession.Close()

	err = createSchema(tempSession)
	if err != nil {
		return nil, err
	}

	cluster.Keyspace = "minerva"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Log.Error("Failed to create session", zap.Error(err))
		return nil, err
	}

	return &Database{Db: session}, nil
}

func createSchema(session *gocql.Session) error {
	err := createKeyspace(session)
	if err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}

	for _, schema := range coreSchema.CoreSchema {
		err = session.Query(schema).Exec()
		if err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	}

	for _, schema := range networkSchema.NetworkSchema {
		err = session.Query(schema).Exec()
		if err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	}

	for _, schema := range contentSchema.ContentSchema {
		err = session.Query(schema).Exec()
		if err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	}

	return nil
}

func createKeyspace(session *gocql.Session) error {
	query := `
	CREATE KEYSPACE IF NOT EXISTS minerva
	WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`

	return session.Query(query).Exec()
}

func (d *Database) Close() {
	d.Db.Close()
}

func (d *Database) GetDatabase() *gocql.Session {
	return d.Db
}
