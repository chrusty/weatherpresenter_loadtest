package resultprocessor

import (
	fmt "fmt"
	http "net/http"
	url "net/url"
	time "time"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) testloadplaylist(result types.Result) types.Result {

	// Prepare the URLs:
	closePlayListURL := fmt.Sprintf("%s/weatherpresenter/ClosePlaylist", result.MachineAddress)
	openPlayListURL := fmt.Sprintf("%s/weatherpresenter/OpenPlaylist", result.MachineAddress)

	// Close the current playlis:
	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress}).Debug("Closing the current playlist")
	tester.HttpClient.Get(closePlayListURL)

	// Open the playlist we're testing:
	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "playlist": tester.Config.PlaylistFile}).Debug("Opening the configured playlist")
	request, _ := http.NewRequest("GET", openPlayListURL, nil)
	query := request.URL.Query()
	query.Add("filepath", tester.Config.PlaylistFile)
	request.URL.RawQuery = query.Encode()
	response, err := tester.HttpClient.Do(request)

	// Calculate how long the request took:
	result.Duration = time.Since(result.RequestTime)

	// // Cast a net/url error (so we can check for timeouts):
	// var netURLError *url.Error

	// Handle the error:
	if netURLError, ok := err.(*url.Error); ok && netURLError.Timeout() {
		result.TimedOut = true
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": err}).Warn("Timed-out making API request to WeatherPresenter")
	} else if err != nil {
		log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "error": err}).Error("Error making API request to WeatherPresenter")
	} else {
		// Append the test result:
		result.ResponseCode = response.StatusCode
	}

	// Fill in the remaining result fields:
	result.Error = err
	result.TestName = "load_playlist"

	return result
}
