package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/models"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type oauth2Controller struct {
	db *sqlx.DB
}

func NewOAuth2Controller(db *sqlx.DB) *oauth2Controller {
	return &oauth2Controller{db: db}
}

func (ctrl oauth2Controller) GetAuth(c *gin.Context) {
	var req GetAuthReq
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("failed to bind request", slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"IsValidConfiguration": false,
		})
		return
	}

	clientName, oauth2ClientScopes, err := ctrl.validateOAuth2ClientConfiguration(c, req.ClientID, req.RedirectURI, req.Scope, req.ResponseType)
	if err != nil {
		slog.Warn("failed to validate oauth2 client configuration", slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"IsValidConfiguration": false,
		})
		return
	}

	c.HTML(http.StatusOK, "auth.tmpl", gin.H{
		"ClientID":             req.ClientID,
		"RedirectURI":          req.RedirectURI,
		"ResponseType":         req.ResponseType,
		"Scope":                req.Scope,
		"State":                req.State,
		"IsValidConfiguration": true,
		"ClientName":           clientName,
		"OAuth2ClientScopes":   oauth2ClientScopes,
	})
}

func (ctrl oauth2Controller) PostAuth(c *gin.Context) {
	var req PostAuthReq
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("failed to bind request", slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"IsValidConfiguration": false,
		})
		return
	}

	clientName, oauth2ClientScopes, err := ctrl.validateOAuth2ClientConfiguration(c, req.ClientID, req.RedirectURI, req.Scope, req.ResponseType)
	if err != nil {
		slog.Warn("failed to validate oauth2 client configuration", slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"IsValidConfiguration": false,
		})
		return
	}

	// if email from input not found
	var foundUser models.User
	if err := ctrl.db.Get(&foundUser, "SELECT id, email, hashed_password FROM users WHERE email=?", req.Email); err != nil {
		slog.Error("user with email not found", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"Error":                "Invalid email or password",
			"Email":                req.Email,
			"ClientID":             req.ClientID,
			"RedirectURI":          req.RedirectURI,
			"ResponseType":         req.ResponseType,
			"Scope":                req.Scope,
			"State":                req.State,
			"IsValidConfiguration": true,
			"ClientName":           clientName,
			"OAuth2ClientScopes":   oauth2ClientScopes,
		})
		return
	}

	// if there is an error while comparing password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(req.Password)); err != nil {
		slog.Error("failed to compare hash and password", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"Error":                "Invalid email or password",
			"Email":                req.Email,
			"ClientID":             req.ClientID,
			"RedirectURI":          req.RedirectURI,
			"ResponseType":         req.ResponseType,
			"Scope":                req.Scope,
			"State":                req.State,
			"IsValidConfiguration": true,
			"ClientName":           clientName,
			"OAuth2ClientScopes":   oauth2ClientScopes,
		})
		return
	}

	authCode, err := ctrl.createAuthCodeWithScopes(c, req.ClientID, foundUser.ID, oauth2ClientScopes)
	if err != nil {
		slog.Error("failed to generate auth code", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"Error":                "Something went wrong",
			"Email":                req.Email,
			"ClientID":             req.ClientID,
			"RedirectURI":          req.RedirectURI,
			"ResponseType":         req.ResponseType,
			"Scope":                req.Scope,
			"State":                req.State,
			"IsValidConfiguration": true,
			"ClientName":           clientName,
			"OAuth2ClientScopes":   oauth2ClientScopes,
		})
		return
	}

	c.Redirect(
		http.StatusTemporaryRedirect,
		fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, authCode, req.State),
	)
}

// Private helpers
// return ClientName, Scopes, error
func (ctrl oauth2Controller) validateOAuth2ClientConfiguration(
	c *gin.Context,
	clientID string,
	redirectURI string,
	scope string,
	responseType string,
) (string, []*models.OAuth2ClientScope, error) {
	if responseType != "code" {
		return "", nil, errors.New("response type should be 'code'")
	}
	var oauth2Client models.OAuth2Client
	if err := ctrl.db.Get(&oauth2Client, "SELECT id, name FROM oauth2_clients WHERE id = ? AND redirect_uri = ?", clientID, redirectURI); err != nil {
		return "", nil, err
	}

	scopeArray := strings.Fields(scope)
	query, args, err := sqlx.In("SELECT * from oauth2_client_scopes WHERE oauth2_client_id = ? AND key_code IN (?)", clientID, scopeArray)
	if err != nil {
		return "", nil, err
	}

	var oauth2ClientScopes []*models.OAuth2ClientScope
	if err := ctrl.db.Select(&oauth2ClientScopes, query, args...); err != nil {
		return "", nil, err
	}

	return oauth2Client.Name, oauth2ClientScopes, nil
}

func (ctrl oauth2Controller) createAuthCodeWithScopes(
	c *gin.Context,
	clientID string,
	userID string,
	scopes []*models.OAuth2ClientScope,
) (string, error) {
	tx, err := ctrl.db.Begin()
	if err != nil {
		return "", nil
	}

	defer tx.Rollback()

	insertAuthCodeQuery := "INSERT INTO oauth2_client_auth_codes(id, code, user_id, oauth2_client_id, expired_at) VALUES(?, ?, ?, ?, ?)"
	authCodeID := ulid.Make().String()
	authCode := ulid.Make().String()
	authCodeExpiredAt := time.Now().Add(time.Hour).Unix()
	if _, err := tx.Exec(insertAuthCodeQuery, authCodeID, authCode, userID, clientID, authCodeExpiredAt); err != nil {
		return "", err
	}

	insertAuthCodeScopesQuery := "INSERT INTO oauth2_client_auth_code_scopes(oauth2_client_auth_code_id, oauth2_client_scope_id) VALUES"
	for idx, scope := range scopes {
		if idx == len(scopes)-1 {
			insertAuthCodeScopesQuery += fmt.Sprintf("('%s','%s');", authCodeID, scope.ID)
		} else {
			insertAuthCodeScopesQuery += fmt.Sprintf("('%s','%s'),", authCodeID, scope.ID)
		}
	}

	if _, err := tx.Exec(insertAuthCodeScopesQuery); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return authCode, nil
}
