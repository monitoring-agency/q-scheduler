package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
)

type UpdateConfigurationRequest struct {
	WorkerPoolCount *uint `json:"worker_pool_count" echotools:"required"`
	ProcessTimeout  *uint `json:"process_timeout" echotools:"required"`
}

func (w *Wrapper) UpdateConfiguration(c echo.Context) error {
	var req UpdateConfigurationRequest
	if err := utility.ValidateJsonForm(c, &req); err != nil {
		return c.String(400, err.Error())
	}

	if *req.WorkerPoolCount == 0 || *req.ProcessTimeout == 0 {
		return c.String(400, "worker_pool_count and process_timeout must be greater than 0")
	}

	w.Configuration.WorkerPoolCount = int(*req.WorkerPoolCount)
	w.Configuration.ProcessTimeout = int(*req.ProcessTimeout)
	w.DB.Save(w.Configuration)
	w.ConfigurationReloadFunc()
	w.Scheduler.Reload()

	return c.NoContent(200)
}
