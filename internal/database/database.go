package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

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

func (db *Database) Select(ctx context.Context, query string, result interface{}, args ...interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a pointer to a slice")
	}

	iter := db.Db.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	// Get the type of the slice element
	sliceElemType := rv.Elem().Type().Elem()

	// Build a mapping of field names to struct fields
	fieldMap := make(map[string]int)
	for i := 0; i < sliceElemType.NumField(); i++ {
		field := sliceElemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			fieldMap[dbTag] = i
		}
	}

	// Iterate over rows returned by the query
	for {
		columns := iter.Columns()
		row := reflect.New(sliceElemType).Elem()
		fieldValues := make([]interface{}, len(columns))

		// Map query columns to struct fields
		for i, column := range columns {
			if fieldIndex, ok := fieldMap[column.Name]; ok {
				field := row.Field(fieldIndex)

				switch field.Type() {
				case reflect.TypeOf(time.Time{}):
					var temp int64
					fieldValues[i] = &temp
				default:
					fieldValues[i] = field.Addr().Interface()
				}
			} else {
				// Placeholder for fields not in the struct
				var temp interface{}
				fieldValues[i] = &temp
			}
		}

		// Scan the row
		if !iter.Scan(fieldValues...) {
			break
		}

		// Convert timestamps if needed
		for i, column := range columns {
			if fieldIndex, ok := fieldMap[column.Name]; ok {
				field := row.Field(fieldIndex)
				if field.Type() == reflect.TypeOf(time.Time{}) {
					intValue := *fieldValues[i].(*int64)
					field.Set(reflect.ValueOf(time.Unix(intValue, 0)))
				}
			}
		}

		rv.Elem().Set(reflect.Append(rv.Elem(), row))
	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}
