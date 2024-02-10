package middlewares

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/constants"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(constants.UserID)
		if userID == nil {
			// Redirect user to login page if couldn't find userID from session
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}

		// Inject userID into context
		c.Set(constants.UserID, userID)
		c.Next()
	}
}
