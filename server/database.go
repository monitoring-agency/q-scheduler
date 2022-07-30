package server

import (
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/myOmikron/q-scheduler/models"
)

var VERSION = "0.1.0"

func initDatabase(config *models.Config) *gorm.DB {
	driver := sqlite.Open(config.Database.Path)
	db := database.Initialize(
		driver,
		models.About{},
		models.Check{},
		models.TimePeriod{},
		models.SchedulingDay{},
		models.SchedulingPeriod{},
	)

	about := models.About{}
	var count int64

	db.Find(&about).Count(&count)
	if count == 0 {
		db.Save(&models.About{
			Version: VERSION,
		})
	} else {
		if about.Version != VERSION {
			color.Println(color.PURPLE, "Updating version")
			about.Version = VERSION
			db.Save(&about)
		}
	}

	return db
}
