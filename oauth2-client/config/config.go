package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     int    `envconfig:"DB_PORT" required:"true"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	Port       int    `envconfig:"PORT" required:"true" default:"8081"`
	// Oauth2 config
	OAuth2ClientID     string `envconfig:"OAUTH2_CLIENT_ID" required:"true"`
	OAuth2ClientSecret string `envconfig:"OAUTH2_CLIENT_SECRET" required:"true"`
	OAuth2AuthURL      string `envconfig:"OAUTH2_AUTH_URL" required:"true"`
	OAuth2TokenURL     string `envconfig:"OAUTH2_TOKEN_URL" required:"true"`
	OAuth2RedirectURL  string `envconfig:"OAUTH2_REDIRECT_URL" required:"true"`
}

func LoadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &c, nil
}
