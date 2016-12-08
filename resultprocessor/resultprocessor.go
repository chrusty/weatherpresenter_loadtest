package resultprocessor

import (
	os "os"
	signal "os/signal"
	sync "sync"
	time "time"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "Result-processor"})
	myConfig config.Config
	report   = Report{
		AverageDuration:    0 * time.Second,
		Errors:             0,
		MaximumDuration:    0 * time.Second,
		MinimumDuration:    24 * time.Hour,
		Results:            []types.Result{},
		Successes:          0,
		TimeOuts:           0,
		TotalDurationNanos: 0,
	}
	reportMutex sync.Mutex
)

func init() {
	// Set up the logrus logger:
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(logrus.DebugLevel)
}

type Report struct {
	AverageDuration    time.Duration
	Errors             int
	MaximumDuration    time.Duration
	MinimumDuration    time.Duration
	Results            []types.Result
	Successes          int
	TimeOuts           int
	TotalDurationNanos int64
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
				if myConfig.Debug {
					log.WithFields(logrus.Fields{"duration": result.Duration, "machine_address": result.MachineAddress, "response_code": result.ResponseCode, "test_name": result.TestName, "timed_out": result.TimedOut}).Debug("Received a result")
				}

				// Lock the report:
				reportMutex.Lock()

				// Add the result to the report:
				report.Results = append(report.Results, result)

				// Update counters and aggregates:
				report.TotalDurationNanos += result.Duration.Nanoseconds()
				report.AverageDuration = time.Duration(report.TotalDurationNanos / int64(len(report.Results)))
				if result.Duration < report.MinimumDuration {
					report.MinimumDuration = result.Duration
				}
				if result.Duration > report.MaximumDuration {
					report.MaximumDuration = result.Duration
				}
				if result.TimedOut {
					report.TimeOuts++
				}
				if result.Error != nil {
					report.Errors++
				} else {
					report.Successes++
				}

				// Unock the report:
				reportMutex.Unlock()

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

func DumpResults() {
	// Get a lock on the report (because it could be changed in the background if we call this when the tests are still running):
	reportMutex.Lock()
	defer reportMutex.Unlock()

	// Log the report:
	log.WithFields(logrus.Fields{"tests_run": len(report.Results)}).Info("Report")
	log.WithFields(logrus.Fields{"test_errors": report.Errors}).Info("Report")
	log.WithFields(logrus.Fields{"test_timeouts": report.TimeOuts}).Info("Report")
	log.WithFields(logrus.Fields{"test_successes": report.Successes}).Info("Report")
	log.WithFields(logrus.Fields{"duration_minimum": report.MinimumDuration}).Info("Report")
	log.WithFields(logrus.Fields{"duration_average": report.AverageDuration}).Info("Report")
	log.WithFields(logrus.Fields{"duration_maximum": report.MaximumDuration}).Info("Report")
}
