package testaburger

import (
	fmt "fmt"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) closePlaylist(result *types.Result) error {

	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress}).Debug("Closing the current playlist")

	// Build a request:
	closePlayListURL := fmt.Sprintf("%s/weatherpresenter/ClosePlaylist", result.MachineAddress)

	// Make the request:
	_, err := tester.HttpClient.Get(closePlayListURL)

	return err
}
