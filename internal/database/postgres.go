package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	networkSchema "github.com/codevault-llc/humblebrag-api/internal/network/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/jmoiron/sqlx"
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
	defer db.Close()

	_, err = db.Exec(turnStringTogether(networkSchema.NetworkSchema))
	if err != nil {
		logger.Log.Error("Failed to create network schema: %v", zap.Error(err))
		return nil, err
	}

	/*err = db.AutoMigrate(&entities.LicenseModel{}, &entities.ScanModel{}, &networkEntities.NetworkModel{},
		&networkEntities.DNSModel{}, &entities.MetadataModel{}, &networkEntities.WhoisModel{}, &contentEntities.FindingModel{},
		&networkEntities.CertificateModel{}, &contentEntities.ContentModel{}, &contentEntities.ContentStorageModel{},
		&contentEntities.ContentTagsModel{}, &contentEntities.ContentAccessLogModel{},
		&entities.FilterModel{}, &entities.RedirectModel{}, &entities.ScreenshotModel{})
	if err != nil {
		logger.Log.Error("Failed to auto migrate entities: %v", zap.Error(err))
		return nil, err
	}*/

	return db, nil
}

// StructToQuery generates an SQL INSERT query from a struct
func StructToQuery(data interface{}, tableName string) (string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the input is a struct
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("input must be a struct")
	}

	// Prepare slices for columns and placeholders
	var columns []string
	var placeholders []string

	// Loop over struct fields to populate columns and placeholders
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		// Get the db tag value; skip fields without a db tag
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" || dbTag == "id" || dbTag == "created_at" || dbTag == "updated_at" || dbTag == "deleted_at" {
			continue
		}

		// Append the column name and placeholder
		columns = append(columns, dbTag)
		placeholders = append(placeholders, ":"+dbTag)
	}

	// Construct the SQL query string
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, nil
}

// InsertStruct inserts a struct into the database and returns the ID
func InsertStruct(tx *sqlx.Tx, query string, data interface{}) (interface{}, error) {
	rows, err := tx.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		return id, nil
	}
	return nil, errors.New("no rows returned")
}
