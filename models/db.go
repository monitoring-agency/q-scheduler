package models

import (
	"github.com/myOmikron/echotools/utilitymodels"
	"time"
)

type About struct {
	utilitymodels.CommonID
	Version string
}

type TimePeriod struct {
	utilitymodels.CommonID
	Start int
	Stop  int
}

type SchedulingDay struct {
	utilitymodels.CommonID
	Day         int
	TimePeriods []*TimePeriod `gorm:"many2many:scheduling_day__time_periods"`
}

type SchedulingPeriod struct {
	utilitymodels.CommonID
	Days []*SchedulingDay `gorm:"many2many:scheduling_period__days"`
}

type Check struct {
	utilitymodels.CommonID
	UUID               string
	Commandline        string
	SchedulingInterval uint
	SchedulingPeriodID uint
	SchedulingPeriod   SchedulingPeriod
}

type Result struct {
	utilitymodels.CommonID
	CheckID       uint
	Check         Check
	Output        string
	ExecutionTime time.Duration
	ReturnCode    int
}
