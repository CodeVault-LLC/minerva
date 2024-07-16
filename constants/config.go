package constants

import "github.com/codevault-llc/humblebrag-api/config"

var VC config.ViperConfig

func InitConfig() {
	err := VC.ReadConfig()

	if err != nil {
		panic(err)
	}
}
