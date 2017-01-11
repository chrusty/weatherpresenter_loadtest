package types

import (
	"fmt"
	"strconv"
	"time"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
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

func (result Result) Print(myConfig config.Config) {
	// Print the result as a CSV line:
	fmt.Printf("\"%s\",\"%s\",\"%s\",%d,%d,%d,%s,\"%s\"\n", result.RequestTime, result.TestName, result.MachineAddress, myConfig.Concurrency, result.Duration, result.ResponseCode, strconv.FormatBool(result.TimedOut), result.Error)
}
