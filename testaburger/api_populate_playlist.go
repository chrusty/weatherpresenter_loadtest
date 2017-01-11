package testaburger

import (
	fmt "fmt"
	http "net/http"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) populatePlaylist(result *types.Result) error {

	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress, "playlist": tester.Config.PlaylistFile}).Debug("Setting presentation-state to 'Editing'")

	// Build a request:
	populatePlayListURL := fmt.Sprintf("http://%s:34567/weatherpresenter/SetPresentationState", result.MachineAddress)
	request, _ := http.NewRequest("GET", populatePlayListURL, nil)

	// Add the query-parameters:
	query := request.URL.Query()
	query.Add("state", "Editing")
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
