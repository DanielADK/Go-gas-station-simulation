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
	AverageDieselTime string
	MaxDieselTime     string
	TotalDieselTime   string
	DieselCount       int
}

type GasStatistics struct {
	AverageGasTime string
	MaxGasTime     string
	TotalGasTime   string
	GasCount       int
}

type LPGStatistics struct {
	AverageLPGTime string
	MaxLPGTime     string
	TotalLPGTime   string
	LPGCount       int
}

type ElectricStatistics struct {
	AverageElectricTime string
	MaxElectricTime     string
	TotalElectricTime   string
	ElectricCount       int
}
type InvalidStatistics struct {
	Count int
}

type RegisterStatistics struct {
	AverageRegisterTime string
	MaxRegisterTime     string
	TotalRegisterTime   string
	RegisterCount       int
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
			AverageDieselTime: fmt.Sprintf("%.2fns", dieselSum/float64(len(s.DieselTimes))),
			MaxDieselTime:     fmt.Sprintf("%.2fns", dieselMax),
			TotalDieselTime:   fmt.Sprintf("%.2fns", dieselSum),
			DieselCount:       len(s.DieselTimes),
		},
		Gas: GasStatistics{
			AverageGasTime: fmt.Sprintf("%.2fns", gasSum/float64(len(s.GasTimes))),
			MaxGasTime:     fmt.Sprintf("%.2fns", gasMax),
			TotalGasTime:   fmt.Sprintf("%.2fns", gasSum),
			GasCount:       len(s.GasTimes),
		},
		LPG: LPGStatistics{
			AverageLPGTime: fmt.Sprintf("%.2fns", lpgSum/float64(len(s.LPGTimes))),
			MaxLPGTime:     fmt.Sprintf("%.2fns", lpgMax),
			TotalLPGTime:   fmt.Sprintf("%.2fns", lpgSum),
			LPGCount:       len(s.LPGTimes),
		},
		Electric: ElectricStatistics{
			AverageElectricTime: fmt.Sprintf("%.2fns", electricSum/float64(len(s.ElectricTimes))),
			MaxElectricTime:     fmt.Sprintf("%.2fns", electricMax),
			TotalElectricTime:   fmt.Sprintf("%.2fns", electricSum),
			ElectricCount:       len(s.ElectricTimes),
		},
		InvalidVehicles: InvalidStatistics{
			Count: s.InvalidVehicles,
		},
		Register: RegisterStatistics{
			AverageRegisterTime: fmt.Sprintf("%.2fns", registerSum/float64(len(s.RegisterTimes))),
			MaxRegisterTime:     fmt.Sprintf("%.2fns", registerMax),
			TotalRegisterTime:   fmt.Sprintf("%.2fns", registerSum),
			RegisterCount:       len(s.RegisterTimes),
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
