package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/scanner"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
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
// @Success 200 {object} responder.APIResponse{data=models.ScanAPIResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans [post]
func CreateScan(w http.ResponseWriter, r *http.Request) {
	license := r.Context().Value("license").(models.LicenseModel)
	if license.ID == 0 {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrAuthInvalidToken))
		return
	}

	var scan models.ScanRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrInvalidRequest))
		return
	}

	if !utils.ValidateURL(scan.Url) {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrInvalidRequest))
		return
	}

	scan.Url = utils.NormalizeURL(scan.Url)
	userAgent := scan.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
	}

	scanResponse, err := scanner.ScanWebsite(scan.Url, userAgent, license.ID)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrScannerFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(models.ConvertScan(scanResponse), "Scan created successfully"))
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
func GetScans(w http.ResponseWriter, r *http.Request) {
	scans, err := service.GetScans()
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(scans, "Scans retrieved successfully"))
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
func GetScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	scan, err := service.GetScan(scanID)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(scan, "Scan retrieved successfully"))
}
