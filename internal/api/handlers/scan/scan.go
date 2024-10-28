package scan

import (
	"encoding/json"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/scanner"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterScanRoutes(router fiber.Router) error {
	router.Get("/scans", GetScans)
	router.Get("/scans/:scanID", GetScan)
	router.Post("/scans", CreateScan)

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
func CreateScan(c *fiber.Ctx) error {
	license := c.Locals("license").(models.LicenseModel)
	if license.ID == 0 {
		return responder.CreateError(responder.ErrAuthInvalidToken).Error
	}

	var scan models.ScanRequest
	err := json.Unmarshal(c.Body(), &scan)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	if !utils.ValidateURL(scan.Url) {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	scan.Url = utils.NormalizeURL(scan.Url)
	userAgent := scan.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
	}

	scanResponse, err := scanner.ScanWebsite(scan.Url, userAgent, license.ID)
	if err != nil {
		return responder.CreateError(responder.ErrScannerFailed).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(models.ConvertScan(scanResponse), "Scan created successfully"))
	return nil
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
	scans, err := service.GetScans()
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(scans, "Scans retrieved successfully"))
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

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	scan, err := service.GetScan(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if scan.ID == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(scan, "Scan retrieved successfully"))
	return nil
}
