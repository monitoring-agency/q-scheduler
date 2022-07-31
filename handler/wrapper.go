package handler

import (
	"gorm.io/gorm"

	"github.com/monitoring-agency/q-scheduler/models"
	"github.com/monitoring-agency/q-scheduler/modules/scheduler"
)

type Wrapper struct {
	Config    *models.Config
	DB        *gorm.DB
	Scheduler scheduler.Scheduler
}
