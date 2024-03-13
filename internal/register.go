package internal

import (
	"hw01/configs"
	"sync"
)

type Register struct {
	Id          int
	TimeRange   [2]float64
	VehicleChan chan *Vehicle
	QueueLen    int
	Mutex       sync.Mutex
}

func InitializeRegisters(configs []configs.RegisterConfig) []*Register {
	var registers []*Register
	for _, cfg := range configs {
		registers = append(registers, &Register{
			Id:          cfg.Id,
			TimeRange:   cfg.TimeRange,
			VehicleChan: make(chan *Vehicle, 10),
			QueueLen:    0,
		})
	}
	return registers
}

func (r *Register) lock() {
	r.Mutex.Lock()
}

func (r *Register) unlock() {
	r.Mutex.Unlock()
}

func (r *Register) Increment() {
	r.lock()
	r.QueueLen++
	r.unlock()
}

func (r *Register) Decrement() {
	r.lock()
	r.QueueLen--
	r.unlock()
}
