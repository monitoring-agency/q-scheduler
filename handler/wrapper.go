package handler

import (
	"gorm.io/gorm"

	"github.com/myOmikron/q-scheduler/models"
	"github.com/myOmikron/q-scheduler/modules/scheduler"
)

type Wrapper struct {
	Config    *models.Config
	DB        *gorm.DB
	Scheduler scheduler.Scheduler
}
