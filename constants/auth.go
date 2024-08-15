package constants

import (
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"os"
)

var (
	DiscordConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}

	DiscordConfigExtension = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URI_EXTENSION"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}
)

func InitAuth() {
	DiscordConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}

	DiscordConfigExtension = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URI_EXTENSION"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}
}
