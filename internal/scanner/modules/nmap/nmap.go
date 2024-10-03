package nmap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func NmapModule(scanId uint, url string) ([]models.PortModel, error) {
	if !utils.IsNmapInstalled() {
		fmt.Println("Nmap is not installed on the system")
		return nil, fmt.Errorf("Nmap is not installed on the system")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	domain := utils.ConvertURLToDomain(url)
	var portStrings []string
	for _, port := range config.PortLists {
		portStrings = append(portStrings, fmt.Sprintf("%d", port.Port))
	}

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(domain),
		nmap.WithPorts(strings.Join(portStrings, ",")),
		nmap.WithMaxRetries(1),
		nmap.WithServiceInfo(),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create Nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("nmap scan timed out after 60 seconds")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to run Nmap scan: %v", err)
	}
	if len(*warnings) > 0 {
		return nil, fmt.Errorf("Nmap scan warnings: %v", warnings)
	}

	var scannedPorts []models.PortModel
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			scannedPorts = append(scannedPorts, models.PortModel{
				Port:  int(port.ID),
				Name:  port.Service.Name,
				State: port.State.String(),
			})
		}
	}

	nmapModel := models.NmapModel{
		ScanID: scanId,
		Hosts:  scannedPorts, // Save ports to the NmapModel
	}

	_, err = service.CreateNmap(nmapModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Nmap model: %v", err)
	}

	return scannedPorts, nil
}
