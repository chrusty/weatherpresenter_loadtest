package resultprocessor

import (
	"os"
	"os/signal"
	"sync"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "Result-processor"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

func Run(myConfigFromMain config.Config, resultsChannel chan types.Result, waitGroup *sync.WaitGroup) {

	log.Info("Starting the result processor")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Handle incoming results:
	go func() {
		for {
			select {

			case result := <-resultsChannel:

				// Log the result:
				log.WithFields(logrus.Fields{"duration": result.Duration, "machine_address": result.MachineAddress, "response_code": result.ResponseCode, "test_name": result.TestName, "timed_out": result.TimedOut}).Debug("Received a result")
			}
		}
	}()

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Warn("Shutting down the result processor")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
