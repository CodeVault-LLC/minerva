package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func SumoLogicAccessID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sumologic-access-id",
		Description: "Discovered a SumoLogic Access ID, potentially compromising log management services and data analytics integrity.",
		// TODO: Make 'su' case-sensitive.
		Regex: generateSemiGenericRegex([]string{"sumo"},
			"su[a-zA-Z0-9]{12}", false),
	}

	// validate
	tps := []string{
		`sumologic.accessId = "su9OL59biWiJu7"`,      // gitleaks:allow
		`sumologic_access_id = "sug5XpdpaoxtOH"`,     // gitleaks:allow
		`export SUMOLOGIC_ACCESSID="suDbJw97o9WVo0"`, // gitleaks:allow
		`SUMO_ACCESS_ID = "suGyI5imvADdvU"`,          // gitleaks:allow
		generateSampleSecret("sumo", "su"+utils.NewSecret(alphaNumeric("12"))),
	}
	fps := []string{
		`- (NSNumber *)sumOfProperty:(NSString *)property;`,
		`- (NSInteger)sumOfValuesInRange:(NSRange)range;`,
		`+ (unsigned char)byteChecksumOfData:(id)arg1;`,
		`.si-sumologic.si--color::before { color: #000099; }`,
		`/// Based on the SumoLogic keyword syntax:`,
		`sumologic_access_id         = ""`,
		`SUMOLOGIC_ACCESSID: ${SUMOLOGIC_ACCESSID}`,
		`export SUMOLOGIC_ACCESSID=XXXXXXXXXXXXXX`, // gitleaks:allow
	}
	return validate(r, tps, fps)
}

func SumoLogicAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sumologic-access-token",
		Description: "Uncovered a SumoLogic Access Token, which could lead to unauthorized access to log data and analytics insights.",
		Regex: generateSemiGenericRegex([]string{"sumo"},
			alphaNumeric("64"), true),

		Keywords: []string{
			"sumo",
		},
	}

	// validate
	tps := []string{
		`export SUMOLOGIC_ACCESSKEY="3HSa1hQfz6BYzlxf7Yb1WKG3Hyovm56LMFChV2y9LgkRipsXCujcLb5ej3oQUJlx"`, // gitleaks:allow
		`SUMO_ACCESS_KEY: gxq3rJQkS6qovOg9UY2Q70iH1jFZx0WBrrsiAYv4XHodogAwTKyLzvFK4neRN8Dk`,             // gitleaks:allow
		`SUMOLOGIC_ACCESSKEY: 9RITWb3I3kAnSyUolcVJq4gwM17JRnQK8ugRaixFfxkdSl8ys17ZtEL3LotESKB7`,         // gitleaks:allow
		`sumo_access_key = "3Kof2VffNQ0QgYIhXUPJosVlCaQKm2hfpWE6F1fT9YGY74blQBIPsrkCcf1TwKE5"`,          // gitleaks:allow
		generateSampleSecret("sumo", utils.NewSecret(alphaNumeric("64"))),
	}
	fps := []string{
		"-e SUMO_ACCESS_KEY=`etcdctl get /sumologic_secret`",
		`SUMO_ACCESS_KEY={SumoAccessKey}`,
		`SUMO_ACCESS_KEY=${SUMO_ACCESS_KEY:=$2}`,
		`sumo_access_key   = "<SUMOLOGIC ACCESS KEY>"`,
		`SUMO_ACCESS_KEY: AbCeFG123`,
	}
	return validate(r, tps, fps)
}
