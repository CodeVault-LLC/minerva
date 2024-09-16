package controller

import (
	"encoding/json"
	"net/http"

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
	router.HandleFunc("/scans", GetScans).Methods("GET")
	router.HandleFunc("/scan", CreateScan).Methods("POST")
}

func CreateScan(w http.ResponseWriter, r *http.Request) {
	var scan models.ScanRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	scan.Url = utils.NormalizeURL(scan.Url)

	scanResponse, err := scanner.ScanWebsite(scan.Url)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan website")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, scanResponse)
}

func GetScans(w http.ResponseWriter, r *http.Request) {
	scans, err := service.GetScans()
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scans")
		return
	}

	utils.RespondWithJSON(w, 200, scans)
}

func GetScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	scan, err := service.GetScan(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan")
		return
	}

	utils.RespondWithJSON(w, 200, scan)
}

func GetScanFindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	findings, err := service.GetScanFindings(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan findings")
		return
	}

	utils.RespondWithJSON(w, 200, findings)
}

func GetScanContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	contents, err := service.GetScanContent(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan contents")
		return
	}

	utils.RespondWithJSON(w, 200, contents)
}
