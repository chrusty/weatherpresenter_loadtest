package types

import (
	"time"
)

type Result struct {
	Duration       time.Duration
	Error          error
	MachineAddress string
	RequestTime    time.Time
	ResponseCode   int
	TestName       string
	TimedOut       bool
}
