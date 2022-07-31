package handler

import "github.com/labstack/echo/v4"

func (w *Wrapper) Reload(c echo.Context) error {
	w.Scheduler.Reload()
	return c.NoContent(200)
}
