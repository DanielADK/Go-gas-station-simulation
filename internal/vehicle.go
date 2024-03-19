package internal

import "sync"

type Vehicle struct {
	Fuel  string
	Done  chan bool
	Mutex sync.Mutex
}

func (s *Vehicle) Lock() {
	s.Mutex.Lock()
}

func (s *Vehicle) Unlock() {
	s.Mutex.Unlock()
}
