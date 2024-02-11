package main

import (
	"log"

	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/config"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/db"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/models"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	mysqlDB, err := db.NewMySqlxDB(conf)
	if err != nil {
		log.Fatalln(err)
	}

	newOauth2Client := &models.OAuth2Client{
		ID:          ulid.Make().String(),
		SecretKey:   ulid.Make().String(),
		Name:        "Some Company",
		RedirectURI: "http://localhost:8081/oauth2/callback",
	}
	insertOauth2ClientQuery := "INSERT INTO oauth2_clients (id, name, secret_key, redirect_uri) VALUES (:id, :name, :secret_key, :redirect_uri)"
	if _, err := mysqlDB.NamedExec(insertOauth2ClientQuery, newOauth2Client); err != nil {
		log.Println("failed to insert oauth2 client")
		log.Fatalln(err)
	}

	insertScopeQuery := "INSERT INTO oauth2_client_scopes (id, key_code, name, description, oauth2_client_id) VALUES (:id, :key_code, :name, :description, :oauth2_client_id);"
	readPhotosScope := &models.OAuth2ClientScope{
		ID:             ulid.Make().String(),
		KeyCode:        "photos.read",
		Name:           "Read photos",
		Description:    "List photos from your account",
		OAuth2ClientID: newOauth2Client.ID,
	}
	if _, err := mysqlDB.NamedExec(insertScopeQuery, readPhotosScope); err != nil {
		log.Println("failed to insert read scope")
		log.Fatalln(err)
	}

	writePhotosScope := &models.OAuth2ClientScope{
		ID:             ulid.Make().String(),
		KeyCode:        "photos.write",
		Name:           "Write photos",
		Description:    "Upload photos to your account",
		OAuth2ClientID: newOauth2Client.ID,
	}
	if _, err := mysqlDB.NamedExec(insertScopeQuery, writePhotosScope); err != nil {
		log.Println("failed to insert write scope")
		log.Fatalln(err)
	}

	password := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password")
		log.Fatalln(err)
	}

	insertUserQuery := "INSERT INTO users (id, email, hashed_password) VALUES (:id, :email, :hashed_password)"
	newUser := &models.User{
		ID:             ulid.Make().String(),
		Email:          "example@mail.com",
		HashedPassword: string(hashedPassword),
	}
	if _, err := mysqlDB.NamedExec(insertUserQuery, newUser); err != nil {
		log.Println("failed to insert user")
		log.Fatalln(err)
	}

	log.Println("completed data seeding")
}
