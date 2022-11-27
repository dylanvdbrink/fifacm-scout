package main

import (
	"fifacm-scout/internal/clients/sqlite"
	"fifacm-scout/internal/controllers"
	"github.com/furkilic/go-boot-config/pkg/go-boot-config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	propErr := gobootconfig.Load()
	if propErr != nil {
		log.Fatalln("Could not retrieve properties")
	}
	sqlite.Initialize()
	router := gin.Default()
	err := router.SetTrustedProxies([]string{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	router.GET("/ping", controllers.Ping)

	router.GET("/player/:id", controllers.GetPlayer)
	router.GET("/players", controllers.GetPlayers)

	err = router.Run()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
