package models

type OAuth2ClientAuthCodeScope struct {
	AuthCodeID string `db:"oauth2_client_auth_code_id"`
	ScopeID    string `db:"oauth2_client_scope_id"`
}
