package config

import (
	_ "embed"

	"github.com/codevault-llc/humblebrag-api/config/lists"
	"github.com/codevault-llc/humblebrag-api/config/rules"
	"github.com/codevault-llc/humblebrag-api/parsers"
	"github.com/codevault-llc/humblebrag-api/types"
	regexp "github.com/wasilibs/go-re2"
)

type ViperConfig struct {
	Description string

	Rules []struct {
		ID          string
		Description string
		Regex       string
		Keywords    []string
	}

	Lists []types.List
}

type Config struct {
	Rules map[string]types.Rule
	Lists map[string]types.List
}

// Order the rules based on alphabetical order of the ID
func (vc *ViperConfig) OrderRules() []types.Rule {
	rules := make([]types.Rule, len(vc.Rules))

	for i, rule := range vc.Rules {
		rules[i] = types.Rule{
			Description: rule.Description,
			RuleID:      rule.ID,
			Regex:       regexp.MustCompile(rule.Regex),
			Keywords:    rule.Keywords,
		}
	}

	return rules
}

var ConfigRules = []*types.Rule{
	rules.AdafruitAPIKey(),
	rules.AdobeClientID(),
	rules.AdobeClientSecret(),
	rules.AgeSecretKey(),
	rules.Airtable(),
	rules.AlgoliaApiKey(),
	rules.AlibabaAccessKey(),
	rules.AlibabaSecretKey(),
	rules.AsanaClientID(),
	rules.AsanaClientSecret(),
	rules.Atlassian(),
	rules.Authress(),
	rules.AWS(),
	rules.BitBucketClientID(),
	rules.BitBucketClientSecret(),
	rules.BittrexAccessKey(),
	rules.BittrexSecretKey(),
	rules.Beamer(),
	rules.CodecovAccessToken(),
	rules.CoinbaseAccessToken(),
	rules.Clojars(),
	rules.CloudflareAPIKey(),
	rules.CloudflareGlobalAPIKey(),
	rules.CloudflareOriginCAKey(),
	rules.ConfluentAccessToken(),
	rules.ConfluentSecretKey(),
	rules.Contentful(),
	rules.Databricks(),
	rules.DatadogtokenAccessToken(),
	rules.DefinedNetworkingAPIToken(),
	rules.DigitalOceanPAT(),
	rules.DigitalOceanOAuthToken(),
	rules.DigitalOceanRefreshToken(),
	rules.DiscordAPIToken(),
	rules.DiscordClientID(),
	rules.DiscordClientSecret(),
	rules.Doppler(),
	rules.DropBoxAPISecret(),
	rules.DropBoxLongLivedAPIToken(),
	rules.DropBoxShortLivedAPIToken(),
	rules.DroneciAccessToken(),
	rules.Duffel(),
	rules.Dynatrace(),
	rules.EasyPost(),
	rules.EasyPostTestAPI(),
	rules.EtsyAccessToken(),
	rules.FacebookSecret(),
	rules.FacebookAccessToken(),
	rules.FacebookPageAccessToken(),
	rules.FastlyAPIToken(),
	rules.FinicityClientSecret(),
	rules.FinicityAPIToken(),
	rules.FlickrAccessToken(),
	rules.FinnhubAccessToken(),
	rules.FlutterwavePublicKey(),
	rules.FlutterwaveSecretKey(),
	rules.FlutterwaveEncKey(),
	rules.FrameIO(),
	rules.FreshbooksAccessToken(),
	rules.GoCardless(),
	// TODO figure out what makes sense for GCP

	rules.GCPServiceAccount(),
	rules.GCPAPIKey(),
	rules.GCPOAuth(),
	rules.GCPGTM(),
	rules.GCPGA(),
	rules.GCPGA2(),

	rules.GitHubPat(),
	rules.GitHubFineGrainedPat(),
	rules.GitHubOauth(),
	rules.GitHubApp(),
	rules.GitHubRefresh(),
	rules.GitlabPat(),
	rules.GitlabPipelineTriggerToken(),
	rules.GitlabRunnerRegistrationToken(),
	rules.GitterAccessToken(),
	rules.GrafanaApiKey(),
	rules.GrafanaCloudApiToken(),
	rules.GrafanaServiceAccountToken(),
	rules.HarnessApiKey(),
	rules.Hashicorp(),
	rules.HashicorpField(),
	rules.Heroku(),
	rules.HubSpot(),
	rules.HuggingFaceAccessToken(),
	rules.HuggingFaceOrganizationApiToken(),
	rules.Intercom(),
	rules.Intra42ClientSecret(),
	rules.JFrogAPIKey(),
	rules.JFrogIdentityToken(),
	rules.JWT(),
	rules.JWTBase64(),
	rules.KrakenAccessToken(),
	rules.KucoinAccessToken(),
	rules.KucoinSecretKey(),
	rules.LaunchDarklyAccessToken(),
	rules.LinearAPIToken(),
	rules.LinearClientSecret(),
	rules.LinkedinClientID(),
	rules.LinkedinClientSecret(),
	rules.LobAPIToken(),
	rules.LobPubAPIToken(),
	rules.MailChimp(),
	rules.MailGunPubAPIToken(),
	rules.MailGunPrivateAPIToken(),
	rules.MailGunSigningKey(),
	rules.MapBox(),
	rules.MattermostAccessToken(),
	rules.MessageBirdAPIToken(),
	rules.MessageBirdClientID(),
	rules.NetlifyAccessToken(),
	rules.NewRelicUserID(),
	rules.NewRelicUserKey(),
	rules.NewRelicBrowserAPIKey(),
	rules.NewRelicInsertKey(),
	rules.NPM(),
	rules.NytimesAccessToken(),
	rules.OktaAccessToken(),
	rules.OpenAI(),
	rules.PlaidAccessID(),
	rules.PlaidSecretKey(),
	rules.PlaidAccessToken(),
	rules.PlanetScalePassword(),
	rules.PlanetScaleAPIToken(),
	rules.PlanetScaleOAuthToken(),
	rules.PostManAPI(),
	rules.Prefect(),
	rules.PrivateKey(),
	rules.PulumiAPIToken(),
	rules.PyPiUploadToken(),
	rules.RapidAPIAccessToken(),
	rules.ReadMe(),
	rules.RubyGemsAPIToken(),
	rules.ScalingoAPIToken(),
	rules.SendbirdAccessID(),
	rules.SendbirdAccessToken(),
	rules.SendGridAPIToken(),
	rules.SendInBlueAPIToken(),
	rules.SentryAccessToken(),
	rules.ShippoAPIToken(),
	rules.ShopifyAccessToken(),
	rules.ShopifyCustomAccessToken(),
	rules.ShopifyPrivateAppAccessToken(),
	rules.ShopifySharedSecret(),
	rules.SidekiqSecret(),
	rules.SidekiqSensitiveUrl(),
	rules.SlackBotToken(),
	rules.SlackUserToken(),
	rules.SlackAppLevelToken(),
	rules.SlackConfigurationToken(),
	rules.SlackConfigurationRefreshToken(),
	rules.SlackLegacyBotToken(),
	rules.SlackLegacyWorkspaceToken(),
	rules.SlackLegacyToken(),
	rules.SlackWebHookUrl(),
	rules.Snyk(),
	rules.StripeAccessToken(),
	rules.SquareAccessToken(),
	rules.SquareSpaceAccessToken(),
	rules.SumoLogicAccessID(),
	rules.SumoLogicAccessToken(),
	rules.TeamsWebhook(),
	rules.TelegramBotToken(),
	rules.TravisCIAccessToken(),
	rules.Twilio(),
	rules.TwitchAPIToken(),
	rules.TwitterAPIKey(),
	rules.TwitterAPISecret(),
	rules.TwitterAccessToken(),
	rules.TwitterAccessSecret(),
	rules.TwitterBearerToken(),
	rules.Typeform(),
	rules.VaultBatchToken(),
	rules.VaultServiceToken(),
	rules.YandexAPIKey(),
	rules.YandexAWSAccessToken(),
	rules.YandexAccessToken(),
	rules.ZendeskSecretKey(),
	rules.GenericCredential(),
	rules.InfracostAPIToken(),

	rules.URLToken(),
	rules.EmailToken(),
}

var ConfigLists = []*types.List{
	{
		Description: "CPBL Filters for ABP & uBO",
		ListID:      "cpbl-abp",
		Categories:  []string{"adblock"},
		Types:       []parsers.ListType{parsers.Domain},
		URL:         "https://raw.githubusercontent.com/bongochong/CombinedPrivacyBlockLists/master/cpbl-abp-list.txt",
		Parser:      lists.CblAbpParser,
	},
	{
		Description: "CPBL Filters for uBO",
		ListID:      "cpbl-ctld",
		Categories:  []string{"adblock"},
		Types:       []parsers.ListType{parsers.Domain},
		URL:         "https://raw.githubusercontent.com/bongochong/CombinedPrivacyBlockLists/master/NoFormatting/cpbl-ctld.txt",
		Parser:      lists.CblCtldParser,
	},
	{
		Description: "A merged hosts file from a variety of other lists.",
		ListID:      "1hosts-pro",
		Categories:  []string{"ads", "crypto", "malware", "privacy"},
		Types:       []parsers.ListType{parsers.Domain},
		URL:         "https://raw.githubusercontent.com/badmojr/1Hosts/master/Pro/hosts.txt",
		Parser:      lists.OneHostsProParser,
	},
	{
		Description: "URLhaus is a project from abuse.ch with the goal of sharing malicious URLs that are being used for malware distribution.",
		ListID:      "urlhaus-abuse-ch",
		Categories:  []string{"malware"},
		Types:       []parsers.ListType{parsers.Domain},
		URL:         "https://urlhaus.abuse.ch/downloads/text/",
		Parser:      lists.URLHausParser,
	},
	{
		Description: "IPsum is a threat intelligence feed based on 30+ different publicly available lists of suspicious and/or malicious IP addresses.",
		ListID:      "ipsum",
		Categories:  []string{"malware"},
		Types:       []parsers.ListType{parsers.IPv4},
		URL:         "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt",
		Parser:      lists.IPSumParser,
	},
}
