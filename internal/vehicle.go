package internal

import "time"

type Vehicle struct {
	Fuel     string
	WaitTime time.Duration
	Done     chan bool
}
