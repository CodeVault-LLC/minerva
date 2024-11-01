package scan

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterModulesRoutes(router fiber.Router) error {
	router.Get("/scans/:scanID/findings", getScanFindings)
	router.Get("/scans/:scanID/contents", getScanContents)
	router.Get("/scans/:scanID/contents/:contentID", getScanContent)
	router.Get("/scans/:scanID/network", getScanNetwork)
	router.Get("/scans/:scanID/metadata", getScanMetadata)

	return nil
}

// @Summary Get scan findings
// @Description Get scan findings
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.FindingResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/findings [get]
func getScanFindings(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	findings, err := repository.FindingRepository.GetScanFindings(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if len(findings) == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertFindings(findings), "Successfully retrieved scan findings"))
	return nil
}

// @Summary Get scan contents
// @Description Get scan contents
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.ContentResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/contents [get]
func getScanContents(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	contents, err := repository.ContentRepository.GetScanContents(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if len(contents) == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(contents, "Successfully retrieved scan contents"))
	return nil
}

// @Summary Get scan content
// @Description Get scan content
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Param contentID path string true "Content ID"
// @Success 200 {object} responder.APIResponse{data=models.ContentResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
func getScanContent(c *fiber.Ctx) error {
	scanID := c.Params("scanID")
	contentID := c.Params("contentID")

	_, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	contentUint, err := utils.ParseUint(contentID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	content, err := repository.ContentRepository.GetScanContent(uint(contentUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if content.ID == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(content, "Successfully retrieved scan content"))
	return nil
}

// @Summary Get scan network
// @Description Get scan network
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.NetworkResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/network [get]
func getScanNetwork(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	network, err := repository.NetworkRepository.GetScanNetwork(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if network.ID == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(network, "Successfully retrieved scan network"))
	return nil
}

// @Summary Get scan metadata
// @Description Get scan metadata
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {object} responder.APIResponse{data=models.MetadataResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/metadata [get]
func getScanMetadata(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	metadata, err := repository.MetadataRepository.GetMetadataByScanID(uint(scanUint))
	if err != nil {
		responder.WriteJSONResponse(c, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if metadata.ID == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertMetadata(metadata), "Successfully retrieved scan metadata"))
	return nil
}
