package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DNSModel struct {
	gorm.Model

	NetworkId uint
	Network   *NetworkModel `gorm:"foreignKey:NetworkId"`

	CNAME       pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ARecords    pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	AAAARecords pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	MXRecords   pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	NSRecords   pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	TXTRecords  pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	PTRRecord   string         `gorm:"not null"`
	DNSSEC      bool           `gorm:"not null"`
}
