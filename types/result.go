package types

import (
	"time"
)

type Result struct {
	Duration       time.Duration
	MachineAddress string
	ResponseCode   int
	TestName       string
	TimedOut       bool
}
