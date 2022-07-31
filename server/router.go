package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/monitoring-agency/q-scheduler/handler"
	"github.com/monitoring-agency/q-scheduler/models"
	"github.com/monitoring-agency/q-scheduler/modules/scheduler"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *models.Config, s scheduler.Scheduler) {
	w := handler.Wrapper{
		Config:      config,
		DB:          db,
		Scheduler:   s,
		ServerStart: time.Now().UTC(),
	}

	e.GET("/api/v1/about", w.About)
	e.POST("/api/v1/reload", w.Reload)
}
