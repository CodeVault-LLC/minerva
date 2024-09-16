package models

import (
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/types"
)

func ConvertUser(user User) UserResponse {
	return UserResponse{
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

func ConvertSubscription(subscription Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                 subscription.ID,
		PlanName:           subscription.PlanName,
		Price:              subscription.Price,
		Currency:           subscription.Currency,
		Interval:           subscription.Interval,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
	}
}

func ConvertUserMinimal(user User) UserMinimalResponse {
	return UserMinimalResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func ConvertScan(scan Scan) ScanAPIResponse {
	return ScanAPIResponse{
		ID:       scan.ID,
		Findings: int64(len(scan.Findings)),

		WebsiteUrl:  scan.WebsiteUrl,
		WebsiteName: scan.WebsiteName,
		Status:      string(scan.Status),
		Sha256:      scan.Sha256,
		SHA1:        scan.SHA1,
		MD5:         scan.MD5,

		Certificates: ConvertCertificates(scan.Certificates),
		Detail:       ConvertDetail(scan.Detail),
		Lists:        ConvertLists(scan.Lists),
		CreatedAt:    scan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    scan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertDetail(detail Detail) DetailResponse {
	return DetailResponse{
		ID:           detail.ID,
		IPAddresses:  detail.IPAddresses,
		IPRanges:     detail.IPRanges,
		DNSNames:     detail.DNSNames,
		PermittedDNS: detail.PermittedDNS,
		ExcludedDNS:  detail.ExcludedDNS,
		HTTPHeaders:  detail.HTTPHeaders,
	}
}

func ConvertList(list List) ListResponse {
	var configList *types.List
	for _, l := range config.ConfigLists {
		if l.ListID == list.ListID {
			configList = l
			break
		}
	}

	return ListResponse{
		ID:          list.ID,
		Description: configList.Description,
		ListID:      list.ListID,
		Categories:  configList.Categories,
		URL:         configList.URL,
	}
}

func ConvertLists(lists []List) []ListResponse {
	var listResponses []ListResponse

	for _, list := range lists {
		listResponses = append(listResponses, ConvertList(list))
	}

	return listResponses
}

func ConvertScans(scans []Scan) []ScanAPIResponse {
	var scanResponses []ScanAPIResponse

	for _, scan := range scans {
		scanResponses = append(scanResponses, ConvertScan(scan))
	}

	return scanResponses
}

func ConvertFindings(findings []Finding) []FindingResponse {
	var findingResponses []FindingResponse

	for _, finding := range findings {
		findingResponses = append(findingResponses, FindingResponse{
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

func ConvertCertificate(certificate Certificate) CertificateResponse {
	return CertificateResponse{
		ID:                 certificate.ID,
		Issuer:             certificate.Issuer,
		Subject:            certificate.Subject,
		NotBefore:          certificate.NotBefore,
		NotAfter:           certificate.NotAfter,
		SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
	}
}

func ConvertCertificates(certificates []Certificate) []CertificateResponse {
	var certificateResponses []CertificateResponse

	for _, certificate := range certificates {
		certificateResponses = append(certificateResponses, CertificateResponse{
			ID:                 certificate.ID,
			Issuer:             certificate.Issuer,
			Subject:            certificate.Subject,
			NotBefore:          certificate.NotBefore,
			NotAfter:           certificate.NotAfter,
			SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
			PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
		})
	}

	if len(certificateResponses) == 0 {
		return []CertificateResponse{}
	}

	return certificateResponses
}

func ConvertContents(content []Content) []ContentResponse {
	var contentResponses []ContentResponse

	for _, c := range content {
		contentResponses = append(contentResponses, ContentResponse{
			ID:      c.ID,
			ScanID:  c.ScanID,
			Name:    c.Name,
			Content: c.Content,
		})
	}

	return contentResponses
}

func ConvertContent(content Content) ContentResponse {
	return ContentResponse{
		ID:      content.ID,
		ScanID:  content.ScanID,
		Name:    content.Name,
		Content: content.Content,
	}
}
