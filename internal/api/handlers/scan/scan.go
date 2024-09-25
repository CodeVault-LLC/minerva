package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/scanner"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gorilla/mux"
)

func RegisterScanRoutes(api *mux.Router) {
	api.HandleFunc("/scans", GetScans).Methods("GET")
	api.HandleFunc("/scans/{scanID}", GetScan).Methods("GET")
	api.HandleFunc("/scans", CreateScan).Methods("POST")
}

// @Summary Create a new scan
// @Description Create a new scan
// @Tags scans
// @Accept json
// @Produce json
// @Param scan body models.ScanRequest true "Scan Request"
// @Success 200 {object} models.ScanAPIResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans [post]
func CreateScan(w http.ResponseWriter, r *http.Request) {
	license := r.Context().Value("license").(models.LicenseModel)
	if license.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var scan models.ScanRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if !utils.ValidateURL(scan.Url) {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	scan.Url = utils.NormalizeURL(scan.Url)

	scanResponse, err := scanner.ScanWebsite(scan.Url, license.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to scan website")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, models.ConvertScan(scanResponse))
}

// @Summary Get all scans
// @Description Get all scans
// @Tags scans
// @Accept json
// @Produce json
// @Success 200 {array} models.ScanAPIResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans [get]
func GetScans(w http.ResponseWriter, r *http.Request) {
	scans, err := service.GetScans()
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scans")
		return
	}

	helper.RespondWithJSON(w, 200, scans)
}

// @Summary Get a scan
// @Description Get a scan
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {object} models.ScanAPIResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans/{scanID} [get]
func GetScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	scan, err := service.GetScan(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan")
		return
	}

	helper.RespondWithJSON(w, 200, scan)
}
