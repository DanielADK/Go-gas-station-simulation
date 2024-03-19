package internal

import (
	"sync"
	"time"
)

type Vehicle struct {
	Fuel             string
	Done             chan bool
	StationQueueNow  time.Time
	RegisterQueueNow time.Time
	Mutex            sync.Mutex
}

func (s *Vehicle) Lock() {
	s.Mutex.Lock()
}

func (s *Vehicle) Unlock() {
	s.Mutex.Unlock()
}
