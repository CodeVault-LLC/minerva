package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/scanner"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gorilla/mux"
)

func RegisterScanRoutes(api *mux.Router) {
	api.HandleFunc("/scans", GetScans).Methods("GET")
	api.HandleFunc("/scans/{scanID}", GetScan).Methods("GET")
	api.HandleFunc("/scan", CreateScan).Methods("POST")
}

func CreateScan(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)

	if !service.CanPerformScan(user, user.Scans) {
		helper.RespondWithError(w, http.StatusForbidden, "Subscription limit reached")
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

	scanResponse, err := scanner.ScanWebsite(scan.Url, user.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to scan website")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, scanResponse)
}

func GetScans(w http.ResponseWriter, r *http.Request) {
	scans, err := service.GetScans()
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scans")
		return
	}

	helper.RespondWithJSON(w, 200, scans)
}

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
