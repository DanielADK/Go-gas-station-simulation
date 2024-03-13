package internal

import (
	"hw01/configs"
	"sync"
)

type Station struct {
	Id          int
	Fuel        string
	TimeRange   [2]float64
	VehicleChan chan *Vehicle
	QueueLen    int
	Mutex       sync.Mutex
}

func InitializeStations(configs []configs.StationConfig) []*Station {
	var stations []*Station
	for _, cfg := range configs {
		stations = append(stations, &Station{
			Id:          cfg.Id,
			Fuel:        cfg.Fuel,
			TimeRange:   cfg.TimeRange,
			VehicleChan: make(chan *Vehicle, 10),
		})
	}
	return stations
}

func (s *Station) lock() {
	s.Mutex.Lock()
}

func (s *Station) unlock() {
	s.Mutex.Unlock()
}

func (s *Station) Increment() {
	s.lock()
	s.QueueLen++
	s.unlock()
}

func (s *Station) Decrement() {
	s.lock()
	s.QueueLen--
	s.unlock()
}
