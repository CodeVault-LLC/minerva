package secrets

import (
	"log"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
)


func ScanSecrets(script []models.ScriptRequest) []utils.RegexReturn {
	var results []utils.RegexReturn

	// Loop through each secret and script, then once finished for the script append the results to the results slice
	for _, rule := range constants.VC.OrderRules() {
		log.Println("Scanning for rule: ", rule.RuleID)
		var scriptResults []utils.Match
		for _, script := range script {
			matches := utils.GenericScan(rule, utils.Script(script))
			if len(matches) > 0 {
				scriptResults = append(scriptResults, matches...)
			}
		}

		if len(scriptResults) > 0 {
			results = append(results, utils.RegexReturn{Name: rule.RuleID, Matches: scriptResults,
			})
		}
	}

	return results
}
