package models

type OAuth2ClientScope struct {
	ID             string `db:"id"`
	KeyCode        string `db:"key_code"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	OAuth2ClientID string `db:"oauth2_client_id"`
}
