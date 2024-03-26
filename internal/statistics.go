package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Statistics struct {
	dieselTime      float64
	dieselMaxTime   float64
	dieselCount     int
	gasTime         float64
	gasMaxTime      float64
	gasCount        int
	lpgTime         float64
	lpgMaxTime      float64
	lpgCount        int
	electricTime    float64
	electricMaxTime float64
	electricCount   int
	registerTime    float64
	registerMaxTime float64
	registerCount   int
	InvalidVehicle  int
	Mutex           sync.Mutex
}

func InitStatistics() *Statistics {
	return &Statistics{
		dieselTime:      0,
		dieselMaxTime:   0,
		dieselCount:     0,
		gasTime:         0,
		gasMaxTime:      0,
		gasCount:        0,
		lpgTime:         0,
		lpgMaxTime:      0,
		lpgCount:        0,
		electricTime:    0,
		electricMaxTime: 0,
		electricCount:   0,
		registerTime:    0,
		registerMaxTime: 0,
		registerCount:   0,
		InvalidVehicle:  0,
	}
}

func (s *Statistics) AddTime(fuel string, time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	switch fuel {
	case "diesel":
		s.dieselCount++
		s.dieselTime += time
		if time > s.dieselMaxTime {
			s.dieselMaxTime = time
		}
		break
	case "gas":
		s.gasCount++
		s.gasTime += time
		if time > s.gasMaxTime {
			s.gasMaxTime = time
		}
		break
	case "lpg":
		s.lpgCount++
		s.lpgTime += time
		if time > s.lpgMaxTime {
			s.lpgMaxTime = time
		}
		break
	case "electric":
		s.electricCount++
		s.electricTime += time
		if time > s.electricMaxTime {
			s.electricMaxTime = time
		}
		break
	}
}

func (s *Statistics) AddInvalidVehicle() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.InvalidVehicle++
}

func (s *Statistics) AddRegisterTime(time float64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.registerTime += time
	if time > s.registerMaxTime {
		s.registerMaxTime = time
	}
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

	return &ExportStatistics{
		Diesel: DieselStatistics{
			AverageDieselQueueTime: fmt.Sprintf("%fs", s.dieselTime/float64(s.dieselCount)),
			MaxDieselQueueTime:     fmt.Sprintf("%fs", s.dieselMaxTime),
			TotalDieselQueueTime:   fmt.Sprintf("%fs", s.dieselTime),
			DieselCount:            s.dieselCount,
		},
		Gas: GasStatistics{
			AverageGasQueueTime: fmt.Sprintf("%fs", s.gasTime/float64(s.gasCount)),
			MaxGasQueueTime:     fmt.Sprintf("%fs", s.gasMaxTime),
			TotalGasQueueTime:   fmt.Sprintf("%fs", s.gasTime),
			GasCount:            s.gasCount,
		},
		LPG: LPGStatistics{
			AverageLPGQueueTime: fmt.Sprintf("%fs", s.lpgTime/float64(s.lpgCount)),
			MaxLPGQueueTime:     fmt.Sprintf("%fs", s.lpgMaxTime),
			TotalLPGQueueTime:   fmt.Sprintf("%fs", s.lpgTime),
			LPGCount:            s.lpgCount,
		},
		Electric: ElectricStatistics{
			AverageElectricQueueTime: fmt.Sprintf("%fs", s.electricTime/float64(s.electricCount)),
			MaxElectricQueueTime:     fmt.Sprintf("%fs", s.electricMaxTime),
			TotalElectricQueueTime:   fmt.Sprintf("%fs", s.electricTime),
			ElectricCount:            s.electricCount,
		},
		InvalidVehicles: InvalidStatistics{
			Count: s.InvalidVehicle,
		},
		Register: RegisterStatistics{
			AverageRegisterQueueTime: fmt.Sprintf("%fs", s.registerTime/float64(s.registerCount)),
			MaxRegisterQueueTime:     fmt.Sprintf("%fs", s.registerMaxTime),
			TotalRegisterQueueTime:   fmt.Sprintf("%fs", s.registerTime),
			RegisterCount:            s.registerCount,
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
