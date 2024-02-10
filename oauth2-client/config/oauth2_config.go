package config

import "golang.org/x/oauth2"

func NewOauth2Config(conf *Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     conf.OAuth2ClientID,
		ClientSecret: conf.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.OAuth2AuthURL,
			TokenURL: conf.OAuth2TokenURL,
		},
		RedirectURL: conf.OAuth2RedirectURL,
		// Depend on scenario, scopes value might be different
		Scopes: []string{"photos.read", "photos.write"},
	}
}
