package handler

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/monitoring-agency/q-scheduler/models"
)

type AboutResponse struct {
	Version      string    `json:"version"`
	LastReloaded time.Time `json:"last_reloaded"`
	Started      time.Time `json:"started"`
}

func (w *Wrapper) About(c echo.Context) error {
	a := models.About{}
	w.DB.Take(&a)
	return c.JSON(200, &AboutResponse{
		Version:      a.Version,
		Started:      w.ServerStart,
		LastReloaded: w.Scheduler.RunningSince(),
	})
}
