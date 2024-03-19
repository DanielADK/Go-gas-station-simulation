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
			VehicleChan: make(chan *Vehicle),
		})
	}
	return stations
}

func (s *Station) Lock() {
	s.Mutex.Lock()
}

func (s *Station) Unlock() {
	s.Mutex.Unlock()
}
