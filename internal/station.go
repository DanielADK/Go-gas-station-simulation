package internal

import (
	"math/rand"
	"time"
)

type Station struct {
	Id        int
	Fuel      string
	TimeRange [2]float64
	Vehicle   chan Vehicle
}

func ServeStation(config StationConfig, vehicles chan Vehicle, toRegister chan Vehicle) {
	for vehicle := range vehicles {
		// Waiting time
		minTime := station.TimeRange[0]
		maxTime := station.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)
		vehicle.WaitTime = time.Second * time.Duration(waitTime)
		// Sleep (fuel loading)
		time.Sleep(vehicle.WaitTime)
		// Add to register queue
		toRegister <- vehicle
		// Wait for register done
		<-vehicle.Done
	}
}
