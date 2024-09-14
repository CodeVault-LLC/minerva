package scan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/scanner/certificate"
	"github.com/codevault-llc/humblebrag-api/scanner/secrets"
	"github.com/codevault-llc/humblebrag-api/scanner/websites"
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

	website, err := websites.ScanWebsite(scan.Url)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, 500, "Failed to scan website")
		return
	}

	secrets := secrets.ScanSecrets(website.Scripts)

	scanModel := models.Scan{
		WebsiteUrl:  scan.Url,
		WebsiteName: website.WebsiteName,

		Sha256: fmt.Sprintf("%x", utils.SHA256(scan.Url)),
		SHA1:   fmt.Sprintf("%x", utils.SHA1(scan.Url)),
		MD5:    fmt.Sprintf("%x", utils.MD5(scan.Url)),

		Status: "pending",
		UserID: user.ID,
	}

	scanResponse, err := controller.CreateScan(scanModel)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to create scan")
		return
	}

	certificate, err := certificate.GetCertificateWebsite(scan.Url, 443)
	if err != nil {
		utils.RespondWithError(w, 500, "Failed to create certificate")
		return
	}

	for _, cert := range certificate {
		err = controller.CreateCertificate(scanResponse.ID, *cert)
		if err != nil {
			utils.RespondWithError(w, 500, "Failed to create certificate")
			return
		}
	}

	controller.CreateFindings(scanResponse.ID, secrets)

	for _, script := range website.Scripts {
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
