package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Statistics struct {
	DieselTimes     []float64
	GasTimes        []float64
	LPGTimes        []float64
	ElectricTimes   []float64
	RegisterTimes   []float64
	InvalidVehicles int
	Mutex           sync.Mutex
}

func InitStatistics() *Statistics {
	return &Statistics{
		DieselTimes:     make([]float64, 0),
		GasTimes:        make([]float64, 0),
		LPGTimes:        make([]float64, 0),
		ElectricTimes:   make([]float64, 0),
		RegisterTimes:   make([]float64, 0),
		InvalidVehicles: 0,
	}
}

func (s *Statistics) AddDieselTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.DieselTimes = append(s.DieselTimes, time)
}

func (s *Statistics) AddGasTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.GasTimes = append(s.GasTimes, time)
}

func (s *Statistics) AddLPGTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.LPGTimes = append(s.LPGTimes, time)
}

func (s *Statistics) AddElectricTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ElectricTimes = append(s.ElectricTimes, time)
}

func (s *Statistics) AddInvalidVehicle() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.InvalidVehicles++
}

func (s *Statistics) AddRegisterTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.RegisterTimes = append(s.RegisterTimes, time)
}

// GetSumAndMaxTime is Mutex-unprotected
func GetSumAndMaxTime(arr []float64) (float64, float64) {
	var sum float64
	var maximum float64
	for _, time := range arr {
		sum += time
		if time > maximum {
			maximum = time
		}
	}
	return sum, maximum
}

type DieselStatistics struct {
	AverageDieselQueueTime string
	MaxDieselQueueTime     string
	TotalDieselQueueTime   string
	DieselCount            int
}

type GasStatistics struct {
	AverageGasQueueTime string
	MaxGasQueueTime     string
	TotalGasQueueTime   string
	GasCount            int
}

type LPGStatistics struct {
	AverageLPGQueueTime string
	MaxLPGQueueTime     string
	TotalLPGQueueTime   string
	LPGCount            int
}

type ElectricStatistics struct {
	AverageElectricQueueTime string
	MaxElectricQueueTime     string
	TotalElectricQueueTime   string
	ElectricCount            int
}
type InvalidStatistics struct {
	Count int
}

type RegisterStatistics struct {
	AverageRegisterQueueTime string
	MaxRegisterQueueTime     string
	TotalRegisterQueueTime   string
	RegisterCount            int
}

type ExportStatistics struct {
	Diesel          DieselStatistics
	Gas             GasStatistics
	LPG             LPGStatistics
	Electric        ElectricStatistics
	Register        RegisterStatistics
	InvalidVehicles InvalidStatistics
}

func (s *Statistics) Export() *ExportStatistics {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	dieselSum, dieselMax := GetSumAndMaxTime(s.DieselTimes)
	gasSum, gasMax := GetSumAndMaxTime(s.GasTimes)
	lpgSum, lpgMax := GetSumAndMaxTime(s.LPGTimes)
	electricSum, electricMax := GetSumAndMaxTime(s.ElectricTimes)
	registerSum, registerMax := GetSumAndMaxTime(s.RegisterTimes)

	return &ExportStatistics{
		Diesel: DieselStatistics{
			AverageDieselQueueTime: fmt.Sprintf("%fs", dieselSum/float64(len(s.DieselTimes))),
			MaxDieselQueueTime:     fmt.Sprintf("%fs", dieselMax),
			TotalDieselQueueTime:   fmt.Sprintf("%fs", dieselSum),
			DieselCount:            len(s.DieselTimes),
		},
		Gas: GasStatistics{
			AverageGasQueueTime: fmt.Sprintf("%fs", gasSum/float64(len(s.GasTimes))),
			MaxGasQueueTime:     fmt.Sprintf("%fs", gasMax),
			TotalGasQueueTime:   fmt.Sprintf("%fs", gasSum),
			GasCount:            len(s.GasTimes),
		},
		LPG: LPGStatistics{
			AverageLPGQueueTime: fmt.Sprintf("%fs", lpgSum/float64(len(s.LPGTimes))),
			MaxLPGQueueTime:     fmt.Sprintf("%fs", lpgMax),
			TotalLPGQueueTime:   fmt.Sprintf("%fs", lpgSum),
			LPGCount:            len(s.LPGTimes),
		},
		Electric: ElectricStatistics{
			AverageElectricQueueTime: fmt.Sprintf("%fs", electricSum/float64(len(s.ElectricTimes))),
			MaxElectricQueueTime:     fmt.Sprintf("%fs", electricMax),
			TotalElectricQueueTime:   fmt.Sprintf("%fs", electricSum),
			ElectricCount:            len(s.ElectricTimes),
		},
		InvalidVehicles: InvalidStatistics{
			Count: s.InvalidVehicles,
		},
		Register: RegisterStatistics{
			AverageRegisterQueueTime: fmt.Sprintf("%fs", registerSum/float64(len(s.RegisterTimes))),
			MaxRegisterQueueTime:     fmt.Sprintf("%fs", registerMax),
			TotalRegisterQueueTime:   fmt.Sprintf("%fs", registerSum),
			RegisterCount:            len(s.RegisterTimes),
		},
	}
}

func (s *ExportStatistics) SaveStatisticsToJSON(filename string) {
	// Object to JSON
	jsonData, err := json.Marshal(*s)
	if err != nil {
		fmt.Println("Chyba při převodu na JSON:", err)
	}

	// Save to file
	err = os.WriteFile(filename, jsonData, 0644)
}
