package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/scanner/secrets"
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
	var user = r.Context().Value("user").(models.User)

	var scan models.ScanRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		utils.RespondWithError(w, 400, "Invalid request")
		return
	}

	secrets := secrets.ScanSecrets(scan.Scripts)

	scanModel := models.Scan{
		WebsiteUrl:  scan.WebsiteUrl,
		WebsiteName: scan.WebsiteName,
		Status:      "pending",
		UserID:      user.ID,
	}

	scanResponse, err := controller.CreateScan(scanModel)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to create scan")
		return
	}

	controller.CreateFindings(scanResponse.ID, secrets)

	for _, script := range scan.Scripts {
		content := models.Content{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := controller.CreateContent(content)
		if err != nil {
			utils.RespondWithError(w, 500, "Failed to create content")
			return
		}
	}

	utils.RespondWithJSON(w, 200, scanResponse)
}

func GetScans(w http.ResponseWriter, r *http.Request) {
	scans, err := controller.GetScans()
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scans")
		return
	}

	utils.RespondWithJSON(w, 200, scans)
}

func GetScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	scan, err := controller.GetScan(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan")
		return
	}

	utils.RespondWithJSON(w, 200, scan)
}

func GetScanFindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	findings, err := controller.GetScanFindings(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan findings")
		return
	}

	utils.RespondWithJSON(w, 200, findings)
}

func GetScanContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	contents, err := controller.GetScanContent(scanID)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to get scan contents")
		return
	}

	utils.RespondWithJSON(w, 200, contents)
}
