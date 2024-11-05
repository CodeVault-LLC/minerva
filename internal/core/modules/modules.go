package modules

import (
	"github.com/codevault-llc/humblebrag-api/internal/core/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
)

type MiniModule interface {
	Run(job entities.JobModel) (interface{}, error) // Executes the mini-module logic
	Name() string                                   // Returns the mini-module name
}

type ScanModule interface {
	Execute(job entities.JobModel, website types.WebsiteAnalysis) error // Executes the module logic
	Name() string                                                       // Returns the module name
}
