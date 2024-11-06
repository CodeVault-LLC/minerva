package repository

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	"github.com/codevault-llc/minerva/internal/contents/models/viewmodels"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
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

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) (uint, error) {
	tx, err := repository.db.Beginx()
	if err != nil {
		return 0, err
	}

	query, values, err := database.StructToQuery(content, "content")
	if err != nil {
		return 0, err
	}

	contentResponse, err := database.InsertStruct(tx, query, values)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return contentResponse, nil
}

func (repository *ContentRepo) FindContentByHash(hashedBody string) (entities.ContentModel, error) {
	query := "SELECT * FROM content WHERE hashed_body = $1"
	stmt, err := repository.db.Preparex(query)
	if err != nil {
		return entities.ContentModel{}, err
	}

	var content entities.ContentModel
	err = stmt.Get(&content, hashedBody)
	if err != nil {
		return entities.ContentModel{}, err
	}

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

func (repository *ContentRepo) CreateContentStorage(storage entities.ContentStorageModel) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	query, values, err := database.StructToQuery(storage, "content_storage")
	if err != nil {
		return err
	}

	_, err = database.InsertStruct(tx, query, values)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *ContentRepo) GetScanContents(scanId uint) ([]viewmodels.Contents, error) {
	var contents []entities.ContentModel
	type CombinedContent struct {
		entities.ContentModel
		entities.ContentStorageModel
	}
	var combinedContents []CombinedContent

	err := repository.db.Select(&combinedContents, "SELECT * FROM content LEFT JOIN content_storage ON content.id = content_storage.content_id WHERE scan_id = $1", scanId)
	if err != nil {
		logger.Log.Error("Failed to retrieve content", zap.Error(err))
		return nil, err
	}

	// Create maps to hold associated tags and storage information.
	contentIDs := make([]uint, len(contents))
	for i, c := range contents {
		contentIDs[i] = c.Id
	}

	// Retrieve associated tags for each content ID.
	tagsMap := make(map[uint][]string)
	var tags []entities.ContentTagsModel
	repository.db.Get(&tags, "SELECT * FROM content_tags WHERE content_id IN $1", contentIDs)

	for _, tag := range tags {
		tagsMap[tag.ContentId] = append(tagsMap[tag.ContentId], tag.Tag)
	}

	storageMap := make(map[uint]entities.ContentStorageModel)
	for _, c := range combinedContents {
		storageMap[c.ContentModel.Id] = c.ContentStorageModel
	}

	contents = make([]entities.ContentModel, len(combinedContents))
	for i, c := range combinedContents {
		contents[i] = c.ContentModel
	}

	// Convert the content models into the content responses with tags and storage details.
	return viewmodels.ConvertContents(contents, tagsMap, storageMap), nil
}

func (repository *ContentRepo) GetScanContent(contentId uint) (entities.ContentModel, error) {
	var content entities.ContentModel

	repository.db.Get(&content, "SELECT * FROM content WHERE id = $1", contentId)

	return content, nil
}
