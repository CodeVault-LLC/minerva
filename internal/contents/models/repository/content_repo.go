package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ContentRepo struct {
	db *sqlx.DB
}

func NewContentRepo(db *sqlx.DB) *ContentRepo {
	return &ContentRepo{db: db}
}

var ContentRepository *ContentRepo

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) (entities.ContentModel, error) {
	tx, err := repository.db.Beginx()
	if err != nil {
		return entities.ContentModel{}, err
	}

	query, err := database.StructToQuery(content, "content")
	if err != nil {
		return entities.ContentModel{}, err
	}

	_, err = database.InsertStruct(tx, query, content)
	if err != nil {
		return entities.ContentModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entities.ContentModel{}, err
	}

	return content, nil
}

func (repository *ContentRepo) FindContentByHash(hashedBody string) (entities.ContentModel, error) {
	var content entities.ContentModel
	repository.db.Get(&content, "SELECT * FROM content WHERE hashed_body = $1", hashedBody)

	return content, nil
}

func (repository *ContentRepo) IncrementAccessCount(contentID uint) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE content SET access_count = access_count + 1 WHERE id = $1", contentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *ContentRepo) AddContentToScan(scanID uint, contentID uint) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO scan_content (scan_id, content_id) VALUES ($1, $2)", scanID, contentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *ContentRepo) CreateContentStorage(storage entities.ContentStorageModel) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	query, err := database.StructToQuery(storage, "content_storage")
	if err != nil {
		return err
	}

	_, err = database.InsertStruct(tx, query, storage)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *ContentRepo) GetScanContents(scanID uint) ([]viewmodels.Contents, error) {
	var content []entities.ContentModel

	// Retrieve the contents for the scan ID.
	repository.db.Get(&content, "SELECT * FROM content WHERE id IN (SELECT content_id FROM scan_content WHERE scan_id = $1)", scanID)

	logger.Log.Info("Retrieved contents for scan", zap.Uint("scanID", scanID), zap.Int("contentCount", len(content)))

	// Create maps to hold associated tags and storage information.
	contentIDs := make([]uint, len(content))
	for i, c := range content {
		contentIDs[i] = c.Id
	}

	// Retrieve associated tags for each content ID.
	tagsMap := make(map[uint][]string)
	var tags []entities.ContentTagsModel
	repository.db.Get(&tags, "SELECT * FROM content_tags WHERE content_id IN $1", contentIDs)

	for _, tag := range tags {
		tagsMap[tag.ContentId] = append(tagsMap[tag.ContentId], tag.Tag)
	}

	// Retrieve associated storage information for each content ID.
	storageMap := make(map[uint]entities.ContentStorageModel)
	var storageRecords []entities.ContentStorageModel
	repository.db.Get(&storageRecords, "SELECT * FROM content_storage WHERE content_id IN $1", contentIDs)

	for _, storageRecord := range storageRecords {
		storageMap[storageRecord.ContentId] = storageRecord
	}

	// Convert the content models into the content responses with tags and storage details.
	return viewmodels.ConvertContents(content, tagsMap, storageMap), nil
}

func (repository *ContentRepo) GetScanContent(contentID uint) (entities.ContentModel, error) {
	var content entities.ContentModel

	// Retrieve the content by ID, preloading the associated tags and storage information.
	repository.db.Get(&content, "SELECT * FROM content WHERE id = $1", contentID)

	return content, nil
}
