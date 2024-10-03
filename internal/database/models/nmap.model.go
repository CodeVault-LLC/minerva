package models

import (
	"gorm.io/gorm"
)

// NmapModel stores the scan results for Nmap.
type NmapModel struct {
	gorm.Model
	ScanID uint
	Scan   *ScanModel `gorm:"foreignKey:ScanID"`

	Hosts []PortModel `gorm:"foreignKey:NmapID"` // Define one-to-many relationship
}

type NmapResponse struct {
	ID uint `json:"id"`

	Hosts []PortResponse `json:"hosts"`
}

// PortModel defines the structure of a scanned port.
type PortModel struct {
	gorm.Model
	NmapID uint   // Foreign key to reference NmapModel
	Port   int    // Port number
	Name   string // Service name
	State  string // Port state (open, closed, etc.)
}

type PortResponse struct {
	ID    uint   `json:"id"`
	Port  int    `json:"port"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func ConvertNmap(nmap NmapModel) NmapResponse {
	var hosts []PortResponse
	for _, host := range nmap.Hosts {
		hosts = append(hosts, PortResponse{
			ID:    host.ID,
			Port:  host.Port,
			Name:  host.Name,
			State: host.State,
		})
	}

	return NmapResponse{
		ID:    nmap.ID,
		Hosts: hosts,
	}
}
