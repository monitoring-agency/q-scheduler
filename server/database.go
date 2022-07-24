package server

import (
	"github.com/myOmikron/echotools/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/myOmikron/q-scheduler/models"
)

func initDatabase(config *models.Config) *gorm.DB {
	driver := sqlite.Open(config.Database.Path)
	db := database.Initialize(driver)
	return db
}
