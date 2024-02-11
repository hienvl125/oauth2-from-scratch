package models

type OAuth2Client struct {
	ID          string `db:"id"`
	SecretKey   string `db:"secret_key"`
	RedirectURI string `db:"redirect_uri"`
	Name        string `db:"name"`
}
