package repository

import (
	"context"

	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	"github.com/codevault-llc/minerva/internal/contents/models/viewmodels"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"go.uber.org/zap"
)

type ContentRepo struct {
	database *database.Database
}

func NewContentRepo(database *database.Database) *ContentRepo {
	return &ContentRepo{database: database}
}

var ContentRepository *ContentRepo

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) error {
	query := "INSERT INTO content (id, hashed_body, scan_id, source, file_size, file_type, storage_type, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	queryResult := repository.database.GetDatabase().Query(query, content.Id, content.HashedBody, content.ScanId, content.Source, content.FileSize, content.FileType, content.StorageType, content.Tags, content.CreatedAt, content.UpdatedAt)
	err := queryResult.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (repository *ContentRepo) FindContentByHash(hashedBody string) (entities.ContentModel, error) {
	ctx := context.Background()
	var contents []entities.ContentModel

	query := "SELECT id FROM minerva.content WHERE hashed_body = ?"
	err := repository.database.Select(ctx, query, &contents, hashedBody)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return entities.ContentModel{}, err
	}

	if len(contents) == 0 {
		return entities.ContentModel{}, nil
	}

	return contents[0], nil
}

func (repository *ContentRepo) CreateContentStorage(storage entities.ContentStorageModel) error {
	query := "INSERT INTO content_storage (id, content_id, object_key, bucket_name, location, storage_endpoint, encryption) VALUES (?, ?, ?, ?, ?, ?, ?)"
	queryResult := repository.database.GetDatabase().Query(query, storage.Id, storage.ContentId, storage.ObjectKey, storage.BucketName, storage.Location, storage.StorageEndpoint, storage.Encryption)
	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to create content storage", zap.Error(err))
		return err
	}

	return nil
}

type CombinedContent struct {
	entities.ContentModel
	entities.ContentStorageModel
}

func (repository *ContentRepo) GetScanContents(scanId string) ([]viewmodels.Contents, error) {
	ctx := context.Background()

	// Struct to hold the combined data fetched in one query
	type CombinedRow struct {
		ContentID       string `cql:"content_id"`
		ScanID          string `cql:"scan_id"`
		ObjectKey       string `cql:"object_key"`
		BucketName      string `cql:"bucket_name"`
		Location        string `cql:"location"`
		StorageEndpoint string `cql:"storage_endpoint"`
		Encryption      string `cql:"encryption"`
	}

	var combinedRows []CombinedRow

	query := `SELECT content.id AS content_id, content.scan_id,
       content_storage.object_key, content_storage.bucket_name,
       content_storage.location, content_storage.storage_endpoint,
       content_storage.encryption FROM content, content_storage WHERE content.scan_id = ? AND content.id = content_storage.content_id ALLOW FILTERING;
`

	err := repository.database.Select(ctx, query, &combinedRows, scanId)
	if err != nil {
		logger.Log.Error("Failed to fetch combined contents", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Fetched Combined Rows", zap.Any("combined_rows", combinedRows))

	// Convert to viewmodel and map as required
	storageMap := make(map[string]entities.ContentStorageModel)
	contents := make([]entities.ContentModel, 0)

	for _, row := range combinedRows {
		contents = append(contents, entities.ContentModel{
			Id:     row.ContentID,
			ScanId: row.ScanID,
		})

		if row.ObjectKey != "" { // Ensure valid storage entry
			storageMap[row.ContentID] = entities.ContentStorageModel{
				ContentId:       row.ContentID,
				ObjectKey:       row.ObjectKey,
				BucketName:      row.BucketName,
				Location:        row.Location,
				StorageEndpoint: row.StorageEndpoint,
				Encryption:      row.Encryption,
			}
		}
	}

	return viewmodels.ConvertContents(contents, storageMap), nil
}

func (repository *ContentRepo) GetScanContent(contentId string) (entities.ContentModel, error) {
	var content entities.ContentModel

	queryResult := repository.database.GetDatabase().Query("SELECT * FROM content WHERE id = ?", contentId)
	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to retrieve content", zap.Error(err))
		return entities.ContentModel{}, err
	}

	return content, nil
}
