package models

type OAuth2ClientAuthCode struct {
	ID             string `db:"id"`
	Code           string `db:"code"`
	UserID         string `db:"user_id"`
	OAuth2ClientID string `db:"oauth2_client_id"`
	ExpiredAt      int64  `db:"expired_at"`
}
