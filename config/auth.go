package config

import (
	"os"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var (
	DiscordConfig = &oauth2.Config{}
)

func InitDiscordAuth() {
	DiscordConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}
}
