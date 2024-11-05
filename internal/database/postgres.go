package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	contentSchema "github.com/codevault-llc/humblebrag-api/internal/contents/models"
	coreSchema "github.com/codevault-llc/humblebrag-api/internal/core/models"
	networkSchema "github.com/codevault-llc/humblebrag-api/internal/network/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func turnStringTogether(strings ...string) string {
	var result string
	for _, str := range strings {
		result += str
	}
	return result
}

func InitPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(turnStringTogether(coreSchema.CoreSchema, networkSchema.NetworkSchema, contentSchema.ContentSchema))
	if err != nil {
		logger.Log.Error("Failed to create network schema: %v", zap.Error(err))
		return nil, err
	}

	return db, nil
}

// StructToQuery generates an SQL INSERT query from a struct
func StructToQuery(data interface{}, tableName string) (string, []interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the input is a struct
	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("input must be a struct")
	}

	// Prepare slices for columns and placeholders
	var columns []string
	var placeholders []string
	var values []interface{}

	// Loop over struct fields to populate columns and placeholders
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		// Get the db tag value; skip fields without a db tag
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" || dbTag == "id" || dbTag == "created_at" || dbTag == "updated_at" || dbTag == "deleted_at" {
			continue
		}

		// Check if the field is a []string type and convert it for PostgreSQL
		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.String {
			// Use pq.Array to wrap the []string value
			values = append(values, pq.Array(fieldValue.Interface()))
		} else {
			// For non-[]string types, add the value directly
			values = append(values, fieldValue.Interface())
		}

		// Append the column name and placeholder
		columns = append(columns, dbTag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}

	// Construct the SQL query string
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values, nil
}

// InsertStruct inserts a struct into the database and returns the ID
func InsertStruct(tx *sqlx.Tx, query string, values []interface{}) (uint, error) {
	// Execute the query with a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Structure to capture the returned ID
	var id struct {
		Val int64 `db:"id"`
	}

	// Execute the query with the provided values
	err := tx.QueryRowxContext(ctx, query, values...).StructScan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot insert into %q: %w", query, err)
	}

	// Log successful insertion
	logger.Log.Info("Inserted record", zap.Int64("id", id.Val))

	return uint(id.Val), nil
}
