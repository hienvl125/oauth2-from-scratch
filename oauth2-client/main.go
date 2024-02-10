package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/config"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/controllers"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/db"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/middlewares"
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

	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("oauth2-client-session", store))

	authController := controllers.NewAuthController(mysqlDB)
	router.GET("/register", authController.GetRegister)
	router.POST("/register", authController.PostRegister)
	router.GET("/login", authController.GetLogin)
	router.POST("/login", authController.PostLogin)

	oauth2Config := config.NewOauth2Config(conf)
	photosController := controllers.NewPhotosController(mysqlDB, oauth2Config)
	router.GET("/photos", middlewares.AuthMiddleware(), photosController.Index)

	router.Run(fmt.Sprintf(":%d", conf.Port))
}
