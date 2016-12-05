package testaburger

import (
	fmt "fmt"

	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

func (tester Tester) triggerPlayout(result *types.Result) error {

	log.WithFields(logrus.Fields{"machine_address": result.MachineAddress}).Debug("Playout: TRIGGER")

	// Build a request:
	triggerPlayoutURL := fmt.Sprintf("%s/weatherpresenter/trigger", result.MachineAddress)

	// Make the request:
	_, err := tester.HttpClient.Get(triggerPlayoutURL)

	return err
}
