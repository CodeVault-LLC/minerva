package scanner

/*
func SaveScan(scan WebsiteScan, userId uint) (models.ScanModel, error) {
	scanModel := models.ScanModel{
		WebsiteUrl:  scan.Website.WebsiteUrl,
		WebsiteName: scan.Website.WebsiteName,

		UserID: userId,

		Sha256: utils.SHA256(scan.Website.WebsiteUrl),
		SHA1:   utils.SHA1(scan.Website.WebsiteUrl),
		MD5:    utils.MD5(scan.Website.WebsiteUrl),

		Status: models.ScanStatusComplete,
	}

	// Create Scan
	scanResponse, err := service.CreateScan(scanModel)
	if err != nil {
		fmt.Println("Failed to create scan", err)
		return models.ScanModel{}, err
	}

	// Create Findings
	service.CreateFindings(scanResponse.ID, scan.Secrets)

	// Create Contents
	for _, script := range scan.Website.Scripts {
		content := models.ContentModel{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Create Networks
	network := models.NetworkModel{
		ScanID:       scanResponse.ID,
		IPAddresses:  scan.IPAddresses,
		HTTPHeaders:  scan.HTTPHeaders,
		IPRanges:     scan.IPRanges,
		DNSNames:     scan.GetDNSScan.CNAME,
		PermittedDNS: scan.GetDNSScan.Permitted,
		ExcludedDNS:  scan.GetDNSScan.Excluded,
	}

	networkResponse, err := service.CreateNetwork(network)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Create Certificates
	for _, certificate := range scan.Certificates {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			fmt.Println("Failed to create certificate", err)
			return models.ScanModel{}, err
		}
	}

	// Create Whois
	whois := models.WhoisModel{
		NetworkId: networkResponse.ID,
		Status: func() string {
			if len(scan.WhoisRecord.Domain.Status) > 0 {
				return scan.WhoisRecord.Domain.Status[0]
			}
			return ""
		}(),

		DomainName:  scan.WhoisRecord.Domain.Name,
		Registrar:   scan.WhoisRecord.Registrar.Name,
		Email:       scan.WhoisRecord.Registrant.Email,
		Phone:       scan.WhoisRecord.Registrant.Phone,
		NameServers: scan.WhoisRecord.Domain.NameServers,

		RegistrantName:       scan.WhoisRecord.Registrant.Name,
		RegistrantCity:       scan.WhoisRecord.Registrant.City,
		RegistrantPostalCode: scan.WhoisRecord.Registrant.PostalCode,
		RegistrantCountry:    scan.WhoisRecord.Registrant.Country,
		RegistrantEmail:      scan.WhoisRecord.Registrant.Email,
		RegistrantPhone:      scan.WhoisRecord.Registrant.Phone,
		RegistrantOrg:        scan.WhoisRecord.Registrant.Organization,

		AdminName:       scan.WhoisRecord.Administrative.Name,
		AdminEmail:      scan.WhoisRecord.Administrative.Email,
		AdminPhone:      scan.WhoisRecord.Administrative.Phone,
		AdminOrg:        scan.WhoisRecord.Administrative.Organization,
		AdminCity:       scan.WhoisRecord.Administrative.City,
		AdminPostalCode: scan.WhoisRecord.Administrative.PostalCode,
		AdminCountry:    scan.WhoisRecord.Administrative.Country,

		Updated: scan.WhoisRecord.Domain.UpdatedDate,
		Created: scan.WhoisRecord.Domain.CreatedDate,
		Expires: scan.WhoisRecord.Domain.ExpirationDate,
	}

	_, err = service.CreateWhois(whois)
	if err != nil {
		return models.ScanModel{}, err
	}

	for _, list := range scan.FoundLists {
		listModel := models.ListModel{
			ScanID: scanResponse.ID,
			ListID: list.ListID,
		}

		_, err := service.CreateList(listModel)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	return scanResponse, nil
}

func UpdateScan(, scanId uint) (models.ScanModel, error) {
	scanModel := models.ScanModel{
		WebsiteUrl:  scan.Website.WebsiteUrl,
		WebsiteName: scan.Website.WebsiteName,

		Sha256: utils.SHA256(scan.Website.WebsiteUrl),
		SHA1:   utils.SHA1(scan.Website.WebsiteUrl),
		MD5:    utils.MD5(scan.Website.WebsiteUrl),

		Status: models.ScanStatusComplete,
	}
	scanModel.ID = scanId

	// Update Scan
	scanResponse, err := service.UpdateScan(scanModel)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Update Findings
	service.UpdateFindings(scanResponse.ID, scan.Secrets)

	// Update Contents
	service.DeleteContents(scanResponse.ID)
	for _, script := range scan.Website.Scripts {
		content := models.ContentModel{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Update Networks
	network := models.NetworkModel{
		ScanID:       scanResponse.ID,
		IPAddresses:  scan.IPAddresses,
		HTTPHeaders:  scan.HTTPHeaders,
		IPRanges:     scan.IPRanges,
		DNSNames:     scan.GetDNSScan.CNAME,
		PermittedDNS: scan.GetDNSScan.Permitted,
		ExcludedDNS:  scan.GetDNSScan.Excluded,
	}

	networkResponse, err := service.UpdateNetwork(network)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Update Certificates
	service.DeleteCertificates(networkResponse.ID)
	for _, certificate := range scan.Certificates {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Create Whois
	whois := models.WhoisModel{
		NetworkId: networkResponse.ID,
		Status: func() string {
			if len(scan.WhoisRecord.Domain.Status) > 0 {
				return scan.WhoisRecord.Domain.Status[0]
			}
			return ""
		}(),

		DomainName:  scan.WhoisRecord.Domain.Name,
		Registrar:   scan.WhoisRecord.Registrar.Name,
		Email:       scan.WhoisRecord.Registrant.Email,
		Phone:       scan.WhoisRecord.Registrant.Phone,
		NameServers: scan.WhoisRecord.Domain.NameServers,

		RegistrantName:       scan.WhoisRecord.Registrant.Name,
		RegistrantCity:       scan.WhoisRecord.Registrant.City,
		RegistrantPostalCode: scan.WhoisRecord.Registrant.PostalCode,
		RegistrantCountry:    scan.WhoisRecord.Registrant.Country,
		RegistrantEmail:      scan.WhoisRecord.Registrant.Email,
		RegistrantPhone:      scan.WhoisRecord.Registrant.Phone,
		RegistrantOrg:        scan.WhoisRecord.Registrant.Organization,

		AdminName:       scan.WhoisRecord.Administrative.Name,
		AdminEmail:      scan.WhoisRecord.Administrative.Email,
		AdminPhone:      scan.WhoisRecord.Administrative.Phone,
		AdminOrg:        scan.WhoisRecord.Administrative.Organization,
		AdminCity:       scan.WhoisRecord.Administrative.City,
		AdminPostalCode: scan.WhoisRecord.Administrative.PostalCode,
		AdminCountry:    scan.WhoisRecord.Administrative.Country,

		Updated: scan.WhoisRecord.Domain.UpdatedDate,
		Created: scan.WhoisRecord.Domain.CreatedDate,
		Expires: scan.WhoisRecord.Domain.ExpirationDate,
	}

	_, err = service.CreateWhois(whois)
	if err != nil {
		return models.ScanModel{}, err
	}

	for _, list := range scan.FoundLists {
		listModel := models.ListModel{
			ScanID: scanResponse.ID,
			ListID: list.ListID,
		}

		_, err := service.CreateList(listModel)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	return scanResponse, nil
}
*/
