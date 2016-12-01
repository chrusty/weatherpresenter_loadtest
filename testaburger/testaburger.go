package resultprocessor

import (
	strings "strings"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "testaburger"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

type Tester struct {
	Config         config.Config
	MachineAddress string
	MachineNumber  int
	ResultsChannel chan types.Result
}

func (tester Tester) Run() {

	// Get the address of the machine we're supposed to test:
	machineAddresses := strings.Split(tester.Config.MachineAddresses, ",")
	if len(machineAddresses) <= tester.MachineNumber {
		log.WithFields(logrus.Fields{"machine_addresses": len(machineAddresses), "machine_number": tester.MachineNumber}).Warn("Not enough machine-addresses provided to support the requested concurrency")
		return
	} else {
		tester.MachineAddress = machineAddresses[tester.MachineNumber]

		log.WithFields(logrus.Fields{"machine_address": tester.MachineAddress, "machine_number": tester.MachineNumber}).Info("Running tests against a WeatherPresenter machine")
	}
}
