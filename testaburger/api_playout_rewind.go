package testaburger

import (
	fmt "fmt"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) rewindPlayout(result *types.Result) error {

	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress}).Debug("Playout: REWIND")

	// Build a request:
	rewindPlayoutURL := fmt.Sprintf("http://%s:34567/weatherpresenter/rewind", result.MachineAddress)

	// Make the request:
	_, err := tester.HttpClient.Get(rewindPlayoutURL)

	return err
}
