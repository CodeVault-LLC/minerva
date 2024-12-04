package repository

import (
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

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) (uint, error) {
	query := "INSERT INTO content (hashed_body, access_count, scan_id) VALUES ($1, $2, $3) RETURNING id"

	queryResult := repository.database.GetDatabase().Query(query, content.HashedBody, content.AccessCount, content.ScanId)
	err := queryResult.Exec()
	if err != nil {
		return 0, err
	}

	var contentId uint
	err = queryResult.Scan(&contentId)
	if err != nil {
		return 0, err
	}

	return contentId, nil
}

func (repository *ContentRepo) FindContentByHash(hashedBody string) (entities.ContentModel, error) {
	query := "SELECT * FROM content WHERE hashed_body = $1"
	queryResult := repository.database.GetDatabase().Query(query, hashedBody)

	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to retrieve content", zap.Error(err))
		return entities.ContentModel{}, err
	}

	var content entities.ContentModel
	err = queryResult.Scan(&content)
	if err != nil {
		logger.Log.Error("Failed to scan content", zap.Error(err))
		return entities.ContentModel{}, err
	}

	return content, nil
}

func (repository *ContentRepo) IncrementAccessCount(contentID uint) error {
	query := "UPDATE content SET access_count = access_count + 1 WHERE id = $1"
	queryResult := repository.database.GetDatabase().Query(query, contentID)

	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to increment access count", zap.Error(err))
		return err
	}

	return nil
}

func (repository *ContentRepo) CreateContentStorage(storage entities.ContentStorageModel) error {
	query := "INSERT INTO content_storage (content_id, storage_key, storage_url) VALUES ($1, $2, $3)"
	queryResult := repository.database.GetDatabase().Query(query)
	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to create content storage", zap.Error(err))
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

	queryResult := repository.database.GetDatabase().Query("SELECT * FROM content WHERE scan_id = $1", scanId)

	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to retrieve content", zap.Error(err))
		return nil, err
	}

	contentIDs := make([]uint, len(contents))
	for i, c := range contents {
		contentIDs[i] = c.Id
	}

	tagsMap := make(map[uint][]string)
	var tags []entities.ContentTagsModel

	queryResult = repository.database.GetDatabase().Query("SELECT * FROM content_tags WHERE content_id IN (?)", contentIDs)
	err = queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to retrieve tags", zap.Error(err))
		return nil, err
	}

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

	queryResult := repository.database.GetDatabase().Query("SELECT * FROM content WHERE id = $1", contentId)
	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to retrieve content", zap.Error(err))
		return entities.ContentModel{}, err
	}

	return content, nil
}
