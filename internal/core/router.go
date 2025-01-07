package core

import (
	"github.com/buger/jsonparser"
	"github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/core/models/repository"
	"github.com/codevault-llc/minerva/internal/core/models/viewmodels"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/responder"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RegisterCoreRouter(router fiber.Router) error {
	router.Get("/scans", GetScans)
	router.Get("/scans/:scanID", GetScan)
	router.Post("/scans", CreateScanHandler(Scheduler))

	router.Get("/jobs/:jobID", GetJob)

	return nil
}

// @Summary Create a new scan
// @Description Create a new scan
// @Tags scans
// @Accept json
// @Produce json
// @Param scan body models.ScanRequest true "Scan Request"
// @Success 200 {object} responder.APIResponse{data=models.ScanAPIResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans [post]
func CreateScanHandler(taskScheduler *TaskScheduler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var scanRequest viewmodels.ScanRequest
		logger.Log.Info("Body", zap.String("body", string(c.Body())))

		if err := jsonparser.ObjectEach(c.Body(), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "url":
				scanRequest.URL = string(value)
			case "user_agent":
				scanRequest.UserAgent = string(value)
			}
			return nil
		}); err != nil {
			return responder.CreateError(responder.ErrInvalidRequest).Error
		}

		if !utils.ValidateURL(scanRequest.URL) {
			return responder.CreateError(responder.ErrInvalidRequest).Error
		}
		scanRequest.URL = utils.NormalizeURL(scanRequest.URL)
		if utils.IsLocalURL(scanRequest.URL) {
			return responder.CreateError(responder.ErrInvalidRequest).Error
		}
		userAgent := scanRequest.UserAgent
		if userAgent == "" {
			userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
		}
		job := entities.JobModel{
			ScanID:    uuid.New().String(),
			ID:        utils.GenerateID(),
			Type:      "WebsiteScan",
			URL:       scanRequest.URL,
			UserAgent: userAgent,
			Status:    entities.Queued,
		}
		taskScheduler.AddJob(&job)
		responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertJob(job), "Scan queued for processing"))
		return nil
	}
}

// @Summary Get all scans
// @Description Get all scans
// @Tags scans
// @Accept json
// @Produce json
// @Success 200 {object} responder.APIResponse{data=[]models.ScanAPIResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans [get]
func GetScans(c *fiber.Ctx) error {
	scans, err := repository.ScanRepository.GetScans()
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if len(scans) == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertScans(scans), "Scans retrieved successfully"))
	return nil
}

// @Summary Get a scan
// @Description Get a scan
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {object} responder.APIResponse{data=models.ScanAPIResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID} [get]
func GetScan(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scan, err := repository.ScanRepository.GetScanResult(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if scan.Id == "" {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertScan(scan), "Scan retrieved successfully"))
	return nil
}

// @Summary Get a job
// @Description Get a job
// @Tags jobs
// @Accept json
// @Produce json
// @Param jobID path string true "Job ID"
// @Success 200 {object} responder.APIResponse{data=models.JobAPIResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /jobs/{jobID} [get]
func GetJob(c *fiber.Ctx) error {
	jobID := c.Params("jobID")

	job := Scheduler.GetArchivedJob(jobID)
	if job == nil {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertJob(*job), "Job retrieved successfully"))
	return nil
}
