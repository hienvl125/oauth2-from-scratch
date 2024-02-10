package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type photosController struct {
	db           *sqlx.DB
	oauth2Config *oauth2.Config
}

func NewPhotosController(db *sqlx.DB, oauth2Config *oauth2.Config) *photosController {
	return &photosController{
		db:           db,
		oauth2Config: oauth2Config,
	}
}

func (ctrl photosController) Index(c *gin.Context) {
	// Actually, the first parameter of AuthCodeURL is an opaque value used by the client to maintain state between the request and callback
	// We need to store it to somewhere and use it to very on callback URL
	OAuth2AuthCodeURL := ctrl.oauth2Config.AuthCodeURL("some-random-state")
	c.HTML(http.StatusOK, "photos_index.tmpl", gin.H{
		"OAuth2AuthCodeURL": OAuth2AuthCodeURL,
	})
}
