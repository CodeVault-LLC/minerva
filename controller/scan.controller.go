package controller

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/scanner"
	"github.com/codevault-llc/humblebrag-api/service"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
)

func ScanRouter(router *mux.Router) {
	router.HandleFunc("/scan/{scanID}", GetScan).Methods("GET")
	router.HandleFunc("/scan/{scanID}/findings", GetScanFindings).Methods("GET")
	router.HandleFunc("/scan/{scanID}/contents", GetScanContents).Methods("GET")
	router.HandleFunc("/scan/{scanID}/network", GetScanNetwork).Methods("GET")
	router.HandleFunc("/scans", GetScans).Methods("GET")
	router.HandleFunc("/scans/statistics", getUserStatisticsHandler).Methods("GET")
	router.HandleFunc("/scan", CreateScan).Methods("POST")
}

func CreateScan(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)

	if !service.CanPerformScan(user.Subscriptions, user.Scans) {
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

func GetScanFindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	findings, err := service.GetScanFindings(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan findings")
		return
	}

	helper.RespondWithJSON(w, 200, findings)
}

func GetScanContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	contents, err := service.GetScanContent(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan contents")
		return
	}

	helper.RespondWithJSON(w, 200, contents)
}

func GetScanNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	network, err := service.GetScanNetwork(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan network")
		return
	}

	helper.RespondWithJSON(w, 200, network)
}

type UserStatisticsResponse struct {
	TotalScans          int64 `json:"totalScans"`
	TotalDomainsScanned int64 `json:"totalDomainsScanned"`
	LastScansIn24Hours  int64 `json:"lastScansIn24Hours"`

	MostScannedDomains []models.ScanAPIResponse `json:"mostScannedWebsites"`
}

func getUserStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)
	if user.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	totalScans, err := service.GetTotalScans()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get total scans")
		return
	}

	totalDomainsScanned, err := service.GetTotalDomainsScanned()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get total domains scanned")
		return
	}

	lastScansIn24Hours, err := service.GetRecentScans()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get last scans in 24 hours")
		return
	}

	mostScannedDomains, err := service.GetMostScannedDomains()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get most scanned domains")
		return
	}

	statistics := UserStatisticsResponse{
		TotalScans:          totalScans,
		TotalDomainsScanned: totalDomainsScanned,
		LastScansIn24Hours:  lastScansIn24Hours,

		MostScannedDomains: models.ConvertScans(mostScannedDomains),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statistics); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to encode statistics response")
	}
}
