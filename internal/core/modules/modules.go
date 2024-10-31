package modules

import "github.com/codevault-llc/humblebrag-api/internal/models/entities"

type MiniModule interface {
	Run(job entities.JobModel) (interface{}, error) // Executes the mini-module logic
	Name() string                                   // Returns the mini-module name
}

type ScanModule interface {
	Execute(job entities.JobModel) error // Executes the module logic
	Name() string                        // Returns the module name
}
