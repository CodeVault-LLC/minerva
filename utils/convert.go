package utils

import "github.com/codevault-llc/humblebrag-api/models"

func ConvertUser(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		// use it as the discord avatar url
		AvatarURL: "https://cdn.discordapp.com/avatars/" + user.DiscordId + "/" + user.Avatar + ".png",
		CreatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertUserMinimal(user models.User) models.UserMinimalResponse {
	return models.UserMinimalResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func ConvertScan(scan models.Scan) models.ScanResponse {
	return models.ScanResponse{
		ID:       scan.ID,
		User:     ConvertUserMinimal(scan.User),
		Findings: int64(len(scan.Findings)),

		WebsiteUrl:  scan.WebsiteUrl,
		WebsiteName: scan.WebsiteName,
		Status:      scan.Status,
		CreatedAt:   scan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   scan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertScans(scans []models.Scan) []models.ScanResponse {
	var scanResponses []models.ScanResponse

	for _, scan := range scans {
		scanResponses = append(scanResponses, ConvertScan(scan))
	}

	return scanResponses
}

func ConvertFindings(findings []models.Finding) []models.FindingResponse {
	var findingResponses []models.FindingResponse

	for _, finding := range findings {
		findingResponses = append(findingResponses, models.FindingResponse{
			ID:     finding.ID,
			ScanID: finding.ScanID,
			Line:   finding.Line,
			Match:  finding.Match,
			Source: finding.Source,

			RegexName:        finding.RegexName,
			RegexDescription: finding.RegexDescription,

			CreatedAt: finding.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: finding.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return findingResponses
}

func ConvertContents(content []models.Content) []models.ContentResponse {
	var contentResponses []models.ContentResponse

	for _, c := range content {
		contentResponses = append(contentResponses, models.ContentResponse{
			ID:      c.ID,
			ScanID:  c.ScanID,
			Name:    c.Name,
			Content: c.Content,
		})
	}

	return contentResponses
}

func ConvertContent(content models.Content) models.ContentResponse {
	return models.ContentResponse{
		ID:      content.ID,
		ScanID:  content.ScanID,
		Name:    content.Name,
		Content: content.Content,
	}
}
