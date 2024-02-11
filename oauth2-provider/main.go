package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/config"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/controllers"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-provider/db"
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

	oauth2Controller := controllers.NewOAuth2Controller(mysqlDB)
	router.GET("/oauth2/auth", oauth2Controller.GetAuth)
	router.POST("/oauth2/auth", oauth2Controller.PostAuth)
	router.Run(fmt.Sprintf(":%d", conf.Port))
}
