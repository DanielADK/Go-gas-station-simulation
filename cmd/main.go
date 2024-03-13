package main

import (
	"fmt"
	"hw01/configs"
	"hw01/internal"
	"math/rand"
	"sync"
	"time"
)

func main() {
	config, err := configs.LoadConfigFromFile("./configs/config.json")
	if err != nil {
		fmt.Println("Chyba při načítání configu:", err)
		return
	}
	fmt.Printf("Config načten s %d pokladnami a %d stanicemi.\n", len(config.Registers), len(config.Stations))

	// Init of stations and registers
	stations := internal.InitializeStations(config.Stations)
	registers := internal.InitializeRegisters(config.Registers)

	// Start vehicle generator
	go generateVehicles(stations, config.Generator)

	var wg sync.WaitGroup

	// run registers
	for _, register := range registers {
		wg.Add(1)
		go registerRoutine(register)
	}
	// run stations
	for _, station := range stations {
		wg.Add(1)
		go stationRoutine(station, registers)
	}

	// waiting till end (never)
	wg.Wait()

	fmt.Println("Simulace čerpací stanice dokončena.")
}

func generateVehicles(stations []*internal.Station, config configs.GeneratorConfig) {
	fuelTypes := config.FuelTypes
	for {
		// Generate vehicle
		chosenFuel := fuelTypes[rand.Intn(len(fuelTypes))]
		vehicle := &internal.Vehicle{
			Fuel: chosenFuel,
			Done: make(chan bool, 1),
		}
		// Choose station
		chosenStation := findShortestQueueStationChan(stations, chosenFuel)
		if chosenStation != nil {
			chosenStation.VehicleChan <- vehicle
			chosenStation.Increment()
			fmt.Printf("[GEN] Vygeneroval jsem vozidlo typu %s a přiřadil jej stanici %d.\n",
				chosenFuel, chosenStation.Id)
		} else {
			fmt.Printf("[GEN] Nemohu najít vhodnou stanici pro vozidlo s palivem %s. Vozidlo nepříjímám\n",
				chosenFuel)
		}

		// Waiting time
		minTime := config.TimeRange[0]
		maxTime := config.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)
		time.Sleep(time.Second * time.Duration(waitTime))
	}
}

func findShortestQueueStationChan(stations []*internal.Station, fuelType string) *internal.Station {
	var minlen int = int(^uint(0) >> 1)
	var chosenStation *internal.Station
	for _, station := range stations {
		if station.Fuel == fuelType && station.QueueLen < minlen {
			chosenStation = station
			minlen = station.QueueLen
		}
	}
	return chosenStation
}

func stationRoutine(station *internal.Station, registerChans []*internal.Register) {
	for vehicle := range station.VehicleChan {
		station.Decrement()
		fmt.Printf("[S%02d] Mám ve frontě %d aut.\n", station.Id, station.QueueLen)

		// Waiting time
		minTime := station.TimeRange[0]
		maxTime := station.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)

		// Print
		fmt.Printf("[S%02d] Vozidlo na %s tankuje po %2.2fs\n", station.Id, vehicle.Fuel, waitTime)

		// Sleep (fuel loading)
		time.Sleep(time.Second * time.Duration(waitTime))

		// Find shortest register queue
		shortestQueueRegisterChan := findShortestQueueRegisterChan(registerChans)
		shortestQueueRegisterChan.VehicleChan <- vehicle
		shortestQueueRegisterChan.Increment()

		// Waiting for payment done
		<-vehicle.Done
		fmt.Printf("[S%02d] Vozidlo %s odjížší.\n", station.Id, vehicle.Fuel)
	}
}

func findShortestQueueRegisterChan(registers []*internal.Register) *internal.Register {
	var minLen int = int(^uint(0) >> 1)
	var chosenRegister *internal.Register

	for _, register := range registers {
		if register.QueueLen < minLen {
			chosenRegister = register
			minLen = register.QueueLen
		}
	}
	return chosenRegister
}

func registerRoutine(register *internal.Register) {
	for vehicle := range register.VehicleChan {
		register.Decrement()
		fmt.Printf("[R%02d] Mám ve frontě %d aut.\n", register.Id, register.QueueLen)

		// Waiting time
		minTime := register.TimeRange[0]
		maxTime := register.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)

		// Print
		fmt.Printf("[R%02d] Vozidlo na %s platí po %2.2fs\n", register.Id, vehicle.Fuel, waitTime)

		// Sleep (payment)
		time.Sleep(time.Second * time.Duration(waitTime))

		// Mark payment done
		vehicle.Done <- true
	}
}
