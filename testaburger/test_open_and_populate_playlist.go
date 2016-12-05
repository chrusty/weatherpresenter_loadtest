package testaburger

import (
	url "net/url"
	time "time"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) testOpenAndPopulatePlaylist(result *types.Result) {

	// Give the result a test-name:
	result.TestName = "open_and_populate_playlist"

	// Close the current playlis:
	tester.closePlaylist(result)

	// Open the configured playlist:
	tester.openPlaylist(result)

	// Cast a net/url error (so we can check for timeouts):
	if netURLError, ok := result.Error.(*url.Error); ok && netURLError.Timeout() {
		result.TimedOut = true
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Warn("Timed-out making API request to WeatherPresenter")
	} else if result.Error != nil {
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Error("Error making API request to WeatherPresenter")
	}

	// Bail if we had any error at all:
	if result.Error != nil {
		// Calculate how long the request took:
		result.Duration = time.Since(result.RequestTime)
		return
	} else {
		// Sleep:
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "sleep_duration": tester.Config.SleepBetweenTests}).Debug("Sleeping between API-calls")
		time.Sleep(tester.Config.SleepBetweenTests)
	}

	// Populate the opened playlist:
	tester.populatePlaylist(result)

	// Cast a net/url error (so we can check for timeouts):
	if netURLError, ok := result.Error.(*url.Error); ok && netURLError.Timeout() {
		result.TimedOut = true
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Warn("Timed-out making API request to WeatherPresenter")
	} else if result.Error != nil {
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Error("Error making API request to WeatherPresenter")
	}

	// Calculate how long the request took:
	result.Duration = time.Since(result.RequestTime) - tester.Config.SleepBetweenTests
}
