package scheduler

import (
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"

	"github.com/monitoring-agency/q-scheduler/models"
)

type Scheduler interface {
	Start()
	Quit()
	Reload()
}

type scheduler struct {
	quit   chan bool
	reload chan bool
	db     *gorm.DB
	pool   worker.Pool
}

func New(db *gorm.DB, config *models.Config) *scheduler {
	return &scheduler{
		quit:   make(chan bool),
		reload: make(chan bool),
		db:     db,
		pool: worker.NewPool(&worker.PoolConfig{
			NumWorker: int(config.Scheduler.Worker),
			QueueSize: int(config.Scheduler.Worker * 10),
		}),
	}
}

func (s *scheduler) loadChecks() []worker.Task {
	var checks []models.Check
	s.db.Preload("SchedulingPeriod.Days.TimePeriods").Find(&checks)
	tasks := make([]worker.Task, 0)

	for _, check := range checks {
		tasks = append(tasks, createCheckTask(check, s))
	}

	return tasks
}

func (s *scheduler) Start() {
	go s.pool.Start()
	s.pool.AddTasks(s.loadChecks())

Loop:
	for {
		select {
		case <-s.reload:
			s.pool.Stop()
			go s.pool.Start()
			s.pool.AddTasks(s.loadChecks())
		case <-s.quit:
			s.pool.Stop()
			break Loop
		}
	}
}

func (s *scheduler) Quit() {
	s.quit <- true
}

func (s *scheduler) Reload() {
	s.reload <- true
}
