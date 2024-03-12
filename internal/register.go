package internal

import (
	"math/rand"
	"time"
)

type Register struct {
	Id        int
	TimeRange [2]float64
	Vehicle   chan Vehicle
}

func ServeRegister(register Register) {
	for vehicle := range register.Vehicle {
		// Waiting time
		minTime := register.TimeRange[0]
		maxTime := register.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)
		// Sleep (payment)
		time.Sleep(time.Second * time.Duration(waitTime))
		// Signal done payment
		vehicle.Done <- true
	}
}
