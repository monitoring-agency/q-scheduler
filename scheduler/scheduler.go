package scheduler

import (
	"gorm.io/gorm"
)

type scheduler struct {
	quit   chan bool
	reload chan bool
}

func New(db *gorm.DB) *scheduler {
	return &scheduler{
		quit:   make(chan bool),
		reload: make(chan bool),
	}
}

func (s *scheduler) Start() {
Loop:
	for {
		select {
		case <-s.quit:
			break Loop
		case <-s.reload:

		}
	}
}

func (s *scheduler) Quit() {
	s.quit <- true
}

func (s *scheduler) Reload() {
	s.reload <- true
}
