package service

import (
	"regexp"

	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"gorm.io/gorm"
)

func CreateScan(scan entities.ScanModel) (entities.ScanModel, error) {
	database.DB.Model(&entities.ScanModel{}).Where("url = ?", scan.Url).Update("status", entities.ScanStatusArchived)

	if err := database.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

// Make this function update the scan information. Now we are just updating some fields of the scan, meaning some can be empty.
func UpdateScan(scan entities.ScanModel) (entities.ScanModel, error) {
	if err := database.DB.Model(&entities.ScanModel{}).Where("id = ?", scan.ID).Updates(scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func GetScans() ([]viewmodels.Scan, error) {
	var scans []entities.ScanModel

	if err := database.DB.Where("status IN (?, ?)", entities.ScanStatusComplete, entities.ScanStatusPending).Order("created_at desc").Find(&scans).Error; err != nil {
		return viewmodels.ConvertScans(scans), err
	}

	return viewmodels.ConvertScans(scans), nil
}

func GetScan(scanID uint) (viewmodels.Scan, error) {
	var scan entities.ScanModel

	if err := database.DB.Where("id = ?", scanID).Preload("Findings").Preload("Lists").Where("status IN (?, ?)", entities.ScanStatusComplete, entities.ScanStatusPending).First(&scan).Error; err != nil {
		return viewmodels.ConvertScan(scan), err
	}

	return viewmodels.ConvertScan(scan), nil
}

func GetScansByUserID(userID uint) ([]entities.ScanModel, error) {
	var scans []entities.ScanModel

	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&scans).
		Error; err != nil {
		return scans, err
	}

	return scans, nil
}

// GetTotalScans returns the total number of scans
func GetTotalScans() (int64, error) {
	var totalScans int64
	if err := database.DB.Model(&entities.ScanModel{}).Count(&totalScans).Error; err != nil {
		return 0, err
	}

	return totalScans, nil
}

// GetTotalDomainsScanned returns the total number of unique domains scanned
func GetTotalDomainsScanned() (int64, error) {
	var totalDomainsScanned int64
	if err := database.DB.Model(&entities.ScanModel{}).Distinct("website_url").Count(&totalDomainsScanned).Error; err != nil {
		return 0, err
	}

	return totalDomainsScanned, nil
}

// GetRecentScans returns the number of scans in the last 24 hours
func GetRecentScans() (int64, error) {
	var lastScansIn24Hours int64
	if err := database.DB.Model(&entities.ScanModel{}).Where("created_at >= ?", utils.Get24HoursAgo()).Count(&lastScansIn24Hours).Error; err != nil {
		return 0, err
	}

	return lastScansIn24Hours, nil
}

// GetMostScannedDomains returns the top 10 most scanned domains
func GetMostScannedDomains() ([]entities.ScanModel, error) {
	var scans []entities.ScanModel
	if err := database.DB.Model(&entities.ScanModel{}).Select("*, count(website_url) as count").Group("website_url, id").Order("count desc").Limit(10).Find(&scans).Error; err != nil {
		return scans, err
	}

	return scans, nil
}

func ExecuteAdvancedQuery(parsedQuery map[string][]string) (interface{}, error) {
	query := database.DB.Model(&entities.ScanModel{})

	for key, values := range parsedQuery {
		switch key {
		case "domain", "websiteurl":
			query = query.Where("website_url LIKE ?", "%"+values[0]+"%")
		case "sha256":
			query = query.Where("sha256 IN ?", values)
		case "sha1":
			query = query.Where("sha1 IN ?", values)
		case "md5":
			query = query.Where("md5 IN ?", values)
		case "ip":
			query = query.Joins("JOIN network_models ON network_models.scan_id = scan_models.id").
				Where("network_models.ip_address IN ?", values)
		case "certificate":
			query = query.Joins("JOIN network_models ON network_models.scan_id = scan_models.id").
				Where("network_models.certificate_sha256 IN ?", values)
		case "status":
			query = query.Where("status IN ?", values)
		case "before":
			query = query.Where("created_at < ?", values[0])
		case "after":
			query = query.Where("created_at > ?", values[0])
		case "default":
			defaultQuery := database.DB.Where("1 = 0") // Start with a false condition
			for _, value := range values {
				defaultQuery = defaultQuery.Or(buildDefaultSearch(value))
			}
			query = query.Where(defaultQuery)
		}
	}

	var results []entities.ScanModel
	if err := query.Preload("Lists").Preload("Findings").Find(&results).Error; err != nil {
		return nil, err
	}

	return viewmodels.ConvertScans(results), nil
}

func buildDefaultSearch(term string) *gorm.DB {
	sha256Regex := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	sha1Regex := regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
	md5Regex := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

	switch {
	case sha256Regex.MatchString(term):
		return database.DB.Where("sha256 = ?", term)
	case sha1Regex.MatchString(term):
		return database.DB.Where("sha1 = ?", term)
	case md5Regex.MatchString(term):
		return database.DB.Where("md5 = ?", term)
	default:
		// If it's not a hash, search in multiple fields
		return database.DB.Where(
			database.DB.Where("website_url LIKE ?", "%"+term+"%").
				Or("website_name LIKE ?", "%"+term+"%").
				Or("sha256 LIKE ?", "%"+term+"%").
				Or("sha1 LIKE ?", "%"+term+"%").
				Or("md5 LIKE ?", "%"+term+"%"),
		)
	}
}
