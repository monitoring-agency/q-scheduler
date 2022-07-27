package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/myOmikron/q-scheduler/handler"
	"github.com/myOmikron/q-scheduler/models"
	"github.com/myOmikron/q-scheduler/modules/scheduler"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *models.Config, s scheduler.Scheduler) {
	w := handler.Wrapper{
		Config:    config,
		DB:        db,
		Scheduler: s,
	}

	e.GET("/api/v1/about", w.About)
}
