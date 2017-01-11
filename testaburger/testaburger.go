package testaburger

import (
	net "net"
	http "net/http"
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
	// Set up the logrus logger:
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
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
	if len(tester.Config.MachineAddresses) <= tester.MachineNumber {
		log.WithFields(logrus.Fields{"machine_addresses": len(tester.Config.MachineAddresses), "machine_number": tester.MachineNumber}).Warn("Not enough machine-addresses provided to support the requested concurrency")
		return
	} else {
		machineAddress := tester.Config.MachineAddresses[tester.MachineNumber]

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
			if tester.Config.TestOpenPopulatePlaylist {
				tester.testOpenAndPopulatePlaylist(&result)
			} else if tester.Config.TestTriggerPlaylist {
				tester.testTriggerPlayout(&result)
			} else {
				log.Warn("You need to choose a test to run! (running Open-Playlist test by default)")
				tester.testOpenPlaylist(&result)
			}

			// Submit the result:
			tester.ResultsChannel <- result

			// Cancel subsequent tests if we got an error:
			if result.Error != nil {
				log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "iteration": iteration}).Warn("Cancelling subsequent tests")
				tester.WaitGroup.Done()
				return
			}

			// Sleep:
			log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "sleep_duration": tester.Config.SleepBetweenTests, "iteration": iteration}).Debug("Sleeping")
			time.Sleep(tester.Config.SleepBetweenTests)

		}

	}

	// Tell main() that we're done:
	tester.WaitGroup.Done()

}
