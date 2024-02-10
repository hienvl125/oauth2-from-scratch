package controllers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/constants"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/models"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type authController struct {
	db *sqlx.DB
}

func NewAuthController(db *sqlx.DB) *authController {
	return &authController{db: db}
}

func (ctrl authController) GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", nil)
}

func (ctrl authController) PostRegister(c *gin.Context) {
	var req PostRegisterReq
	if err := c.Bind(&req); err != nil {
		slog.Error("failed to bind request", slog.Any("error", err))
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Invalid parameters",
		})
		return
	}

	// Validate password and password confirmation inputs
	if req.Password != req.PasswordConfirmation {
		slog.Error("doesn't match password confirmation")
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Password doesn't match password confirmation",
			"Email": req.Email,
		})
		return
	}

	// Check a user existed or not?
	var foundUser models.User
	err := ctrl.db.Get(&foundUser, "SELECT id FROM users WHERE email=?", req.Email)
	if err == nil {
		slog.Error("user with email already existed", slog.String("email", req.Email))
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Email already be taken",
			"Email": req.Email,
		})
		return
	}

	// If error is not 'no rows' error, it might be other database related error?
	if err != sql.ErrNoRows {
		slog.Error("failed to query user by email", slog.String("email", req.Email))
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Something went wrong",
			"Email": req.Email,
		})
		return
	}

	// Hash password from password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash user's password", slog.Any("error", err))
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Something went wrong",
			"Email": req.Email,
		})
	}

	newUser := models.User{
		ID:             ulid.Make().String(),
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	}
	insertUserQuery := "INSERT INTO users (id, email, hashed_password) VALUES (:id, :email, :hashed_password)"
	if _, insertUserErr := ctrl.db.NamedExec(insertUserQuery, newUser); insertUserErr != nil {
		slog.Error("failed to insert user into database", slog.Any("error", insertUserErr), slog.String("email", req.Email))
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"Error": "Something went wrong",
			"Email": req.Email,
		})
		return
	}

	// Redirect to login page if registered account successfully
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func (ctrl authController) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func (ctrl authController) PostLogin(c *gin.Context) {
	var req PostLoginReq
	if err := c.Bind(&req); err != nil {
		slog.Error("failed to bind request", slog.Any("error", err))
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Error": "Invalid parameters",
		})
		return
	}

	// if email from input not found
	var foundUser models.User
	if err := ctrl.db.Get(&foundUser, "SELECT id, email, hashed_password FROM users WHERE email=?", req.Email); err != nil {
		slog.Error("user with email not found", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Error": "Invalid email or password",
			"Email": req.Email,
		})
		return
	}

	// if there is an error while comparing password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(req.Password)); err != nil {
		slog.Error("failed to compare hash and password", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Error": "Invalid email or password",
			"Email": req.Email,
		})
		return
	}

	// store authenticated user id into session
	session := sessions.Default(c)
	session.Set(constants.UserID, foundUser.ID)
	if err := session.Save(); err != nil {
		slog.Error("failed to store authenticated user into session", slog.String("email", req.Email), slog.Any("error", err))
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Error": "Something went wrong",
			"Email": req.Email,
		})
		return
	}

	// redirect to photos page
	c.Redirect(http.StatusMovedPermanently, "/photos")
}
