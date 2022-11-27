package sqlite

import (
	"fifacm-scout/internal/models"
	gobootconfig "github.com/furkilic/go-boot-config/pkg/go-boot-config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Initialize() {
	dbLocation := gobootconfig.GetStringWithDefault("db_path", "fifacm.db")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	log.Println("running automigrate")
	migrateErr := DB.AutoMigrate(&models.Player{}, &models.DBUpdate{})
	if migrateErr != nil {
		panic("failed to migrate: " + migrateErr.Error())
	}
}
