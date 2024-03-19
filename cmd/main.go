package main

import (
	"fmt"
	"hw01/configs"
	"hw01/internal"
	"math/rand"
	"sync"
	"time"
)

type GasStation struct {
	Stations  []*internal.Station
	Registers []*internal.Register
	done      chan bool
	wg        sync.WaitGroup
	finishWg  sync.WaitGroup
}

func main() {
	config, err := configs.LoadConfigFromFile("./configs/config-test.json")
	if err != nil {
		fmt.Println("Chyba při načítání configu:", err)
		return
	}
	fmt.Printf("Config načten s %d pokladnami a %d stanicemi.\n", len(config.Registers), len(config.Stations))

	// Init of stations and registers
	var gasStation = GasStation{
		Stations:  internal.InitializeStations(config.Stations),
		Registers: internal.InitializeRegisters(config.Registers),
		done:      make(chan bool),
	}

	// run registers
	for _, register := range gasStation.Registers {
		gasStation.wg.Add(1)
		gasStation.finishWg.Add(1)
		go func(register *internal.Register) {
			registerRoutine(register, &gasStation)
		}(register)
	}
	// run stations
	for _, station := range gasStation.Stations {
		gasStation.wg.Add(1)
		gasStation.finishWg.Add(1)
		go func(station *internal.Station) {
			stationRoutine(station, &gasStation)
		}(station)
	}

	// Start vehicle generator
	generateVehicles(&gasStation, config.Generator)

	// K dogenerování
	gasStation.wg.Wait()

	// Ukončení simulace
	gasStation.finishWg.Wait()

	fmt.Println("Simulace čerpací stanice dokončena.")
	// run registers
	for _, register := range gasStation.Registers {
		fmt.Printf("Pokladna %d končí s %d vozidly.\n", register.Id, register.QueueLen)
	}
	// run stations
	for _, station := range gasStation.Stations {
		fmt.Printf("Stanice %d končí s %d vozidly.\n", station.Id, station.QueueLen)
	}
}

func generateVehicles(gasStation *GasStation, config configs.GeneratorConfig) {
	// Generate vehicles
	fuelTypes := config.FuelTypes
	var i uint
	for i = 0; i < config.CountOfCars; i++ {
		fmt.Printf("[GEN] Generuji vozidlo %06d/%06d\n", i+1, config.CountOfCars)
		// Generate vehicle
		chosenFuel := fuelTypes[rand.Intn(len(fuelTypes))]
		vehicle := &internal.Vehicle{
			Fuel: chosenFuel,
			Done: make(chan bool, 1),
		}
		// Choose station
		chosenStation := findShortestQueueStationChan(gasStation.Stations, chosenFuel)
		if chosenStation != nil {
			// Send vehicle to station
			chosenStation.VehicleChan <- vehicle
			chosenStation.Lock()
			chosenStation.QueueLen++
			chosenStation.Unlock()

			//fmt.Printf("[GEN] Vygeneroval jsem vozidlo typu %s a přiřadil jej stanici %d.\n", chosenFuel, chosenStation.Id)
		} else {
			fmt.Printf("[GEN] Nemohu najít vhodnou stanici pro vozidlo s palivem %s. Vozidlo nepříjímám\n", chosenFuel)
		}

		// Waiting time
		minTime := config.TimeRange[0]
		maxTime := config.TimeRange[1]
		waitTime := minTime + rand.Float64()*(maxTime-minTime)
		time.Sleep(time.Second * time.Duration(waitTime))
	}
	close(gasStation.done)
}

func areAllStationsEmptyQueue(stations []*internal.Station) bool {
	lockAllStations(stations)
	for _, station := range stations {
		if station.QueueLen > 0 {
			unlockAllStations(stations)
			return false
		}
	}
	unlockAllStations(stations)
	return true
}

func lockAllRegisters(registers []*internal.Register) {
	for _, register := range registers {
		register.Lock()
	}
}

func unlockAllRegisters(registers []*internal.Register) {
	for _, register := range registers {
		register.Unlock()
	}
}

func lockAllStations(stations []*internal.Station) {
	for _, station := range stations {
		station.Lock()
	}
}

func unlockAllStations(stations []*internal.Station) {
	for _, station := range stations {
		station.Unlock()
	}
}

func findShortestQueueStationChan(stations []*internal.Station, fuelType string) *internal.Station {
	lockAllStations(stations)

	var minlen int
	var chosenStation *internal.Station = nil

	for _, station := range stations {
		if station.Fuel == fuelType && (station.QueueLen < minlen || chosenStation == nil) {
			chosenStation = station
			minlen = station.QueueLen
		}
	}
	unlockAllStations(stations)
	return chosenStation
}

func stationRoutine(station *internal.Station, gasStation *GasStation) {
	fmt.Printf("[S%02d] Otevírám stanici.\n", station.Id)
	gasStation.wg.Done()
	for {
		select {
		case vehicle := <-station.VehicleChan:
			/*
				station.Lock()
				fmt.Printf("[S%02d] Mám ve frontě %d aut.\n", station.Id, station.QueueLen)
				station.Unlock()
			*/

			// Waiting time
			minTime := station.TimeRange[0]
			maxTime := station.TimeRange[1]
			waitTime := minTime + rand.Float64()*(maxTime-minTime)

			// Print
			//fmt.Printf("[S%02d] Vozidlo na %s tankuje po %2.2fs\n", station.Id, vehicle.Fuel, waitTime)

			// Sleep (fuel loading)
			time.Sleep(time.Second * time.Duration(waitTime))

			// Find shortest register queue
			shortestQueueRegisterChan := findShortestQueueRegisterChan(gasStation.Registers)
			// Send vehicle to register
			shortestQueueRegisterChan.VehicleChan <- vehicle
			shortestQueueRegisterChan.Lock()
			shortestQueueRegisterChan.QueueLen++
			shortestQueueRegisterChan.Unlock()

			// Waiting for payment done
			<-vehicle.Done

			// Decrement queue length
			station.Lock()
			station.QueueLen--
			station.Unlock()
			//fmt.Printf("[S%02d] Vozidlo %s odjížší.\n", station.Id, vehicle.Fuel)
		case <-gasStation.done:
			station.Lock()
			if station.QueueLen == 0 {
				station.Unlock()
				fmt.Printf("[S%02d] Stanice končí.\n", station.Id)
				close(station.VehicleChan)
				gasStation.finishWg.Done()
				return
			}
			station.Unlock()
		}
	}
}

func findShortestQueueRegisterChan(registers []*internal.Register) *internal.Register {
	lockAllRegisters(registers)

	var minLen int
	var chosenRegister *internal.Register = nil

	for _, register := range registers {
		if register.QueueLen < minLen || chosenRegister == nil {
			chosenRegister = register
			minLen = register.QueueLen
		}
	}
	unlockAllRegisters(registers)
	return chosenRegister
}

func registerRoutine(register *internal.Register, gasStation *GasStation) {
	fmt.Printf("[R%02d] Otevírám pokladnu.\n", register.Id)
	gasStation.wg.Done()
	for {
		select {
		case vehicle := <-register.VehicleChan:
			/*
				register.Lock()
				fmt.Printf("[R%02d] Mám ve frontě %d aut.\n", register.Id, register.QueueLen)
				register.Unlock()
			*/

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

			// Decrement queue length
			register.Lock()
			register.QueueLen--
			register.Unlock()
		case <-gasStation.done:
			register.Lock()
			if register.QueueLen == 0 && areAllStationsEmptyQueue(gasStation.Stations) {
				register.Unlock()
				fmt.Printf("[R%02d] Pokladna končí.\n", register.Id)
				close(register.VehicleChan)
				gasStation.finishWg.Done()
				return
			}
			register.Unlock()
		}
	}
}
