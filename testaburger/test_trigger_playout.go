package testaburger

import (
	url "net/url"
	time "time"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) testTriggerPlayout(result *types.Result) {

	// Give the result a test-name:
	result.TestName = "trigger_playout"

	// Rewind playout on the current playlist:
	tester.rewindPlayout(result)

	// Trigger playout on the current playlist:
	tester.triggerPlayout(result)

	// Calculate how long the request took:
	result.Duration = time.Since(result.RequestTime)

	// Cast a net/url error (so we can check for timeouts):
	if netURLError, ok := result.Error.(*url.Error); ok && netURLError.Timeout() {
		result.TimedOut = true
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Warn("Timed-out making API request to WeatherPresenter")
	} else if result.Error != nil {
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": result.Error, "stage": result.TestName}).Error("Error making API request to WeatherPresenter")
	}

}
