package testaburger

import (
	fmt "fmt"
	http "net/http"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) openPlaylist(result *types.Result) error {

	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "playlist": tester.Config.PlaylistFile}).Debug("Opening the configured playlist")

	// Build a request:
	openPlayListURL := fmt.Sprintf("%s/weatherpresenter/OpenPlaylist", result.MachineAddress)
	request, _ := http.NewRequest("GET", openPlayListURL, nil)

	// Add the query-parameters:
	query := request.URL.Query()
	query.Add("filepath", tester.Config.PlaylistFile)
	request.URL.RawQuery = query.Encode()

	// Make the request & append the results:
	response, err := tester.HttpClient.Do(request)
	if err == nil {
		result.ResponseCode = response.StatusCode
	}

	// Set the results error:
	result.Error = err

	return err
}
