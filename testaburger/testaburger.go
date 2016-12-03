package testaburger

import (
	net "net"
	http "net/http"
	strings "strings"
	sync "sync"
	time "time"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"logger": "testaburger"})
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

type Tester struct {
	Config         config.Config
	HttpClient     *http.Client
	MachineNumber  int
	ResultsChannel chan types.Result
	WaitGroup      *sync.WaitGroup
}

func (tester Tester) Run() {

	// Get the address of the machine we're supposed to test:
	machineAddresses := strings.Split(tester.Config.MachineAddresses, ",")
	if len(machineAddresses) <= tester.MachineNumber {
		log.WithFields(logrus.Fields{"machine_addresses": len(machineAddresses), "machine_number": tester.MachineNumber}).Warn("Not enough machine-addresses provided to support the requested concurrency")
		return
	} else {
		machineAddress := machineAddresses[tester.MachineNumber]

		log.WithFields(logrus.Fields{"machine_address": machineAddress, "machine_number": tester.MachineNumber}).Info("Running tests against a WeatherPresenter machine")

		// Make an HTTP client with the configured timeout and keepalive:
		tester.HttpClient = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   tester.Config.APITimeout,
					KeepAlive: tester.Config.APIKeepAlive,
				}).Dial,
				TLSHandshakeTimeout:   tester.Config.APITimeout,
				ResponseHeaderTimeout: tester.Config.APITimeout,
				ExpectContinueTimeout: tester.Config.APITimeout,
			},
		}

		// Run the tests:
		for iteration := 0; iteration < tester.Config.Iterations; iteration++ {

			// Prepare a test-result:
			result := types.Result{
				MachineAddress: machineAddress,
				RequestTime:    time.Now(),
			}

			// Run the appropriate test:
			if tester.Config.TestOpenPlaylist {
				tester.testOpenPlaylist(&result)
			} else if tester.Config.TestOpenPopulatePlaylist {
				tester.testOpenAndPopulatePlaylist(&result)
			} else {
				log.Error("You need to choose a test to run!")
			}

			// Submit the result:
			tester.ResultsChannel <- result

			// Sleep:
			log.WithFields(logrus.Fields{"machine_addresses": len(machineAddresses), "sleep_duration": tester.Config.SleepBetweenTests, "iteration": iteration}).Debug("Sleeping")
			time.Sleep(tester.Config.SleepBetweenTests)

		}

	}

	// Tell main() that we're done:
	tester.WaitGroup.Done()

}
