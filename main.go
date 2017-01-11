package main

import (
	flag "flag"
	sync "sync"
	// time "time"

	config "github.com/chrusty/weatherpresenter_loadtest/config"
	resultprocessor "github.com/chrusty/weatherpresenter_loadtest/resultprocessor"
	testaburger "github.com/chrusty/weatherpresenter_loadtest/testaburger"
	types "github.com/chrusty/weatherpresenter_loadtest/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	conf      config.Config
	log       = logrus.WithFields(logrus.Fields{"logger": "main"})
	waitGroup = &sync.WaitGroup{}
	testers   = make([]testaburger.Tester, 100)
)

func init() {
	// Process the command-line parameters:
	flag.StringVar(&conf.ConfFile, "configfile", "loadtest.yaml", "Name of a config-file to load")
	flag.BoolVar(&conf.Debug, "debug", false, "Run in DEBUG mode")
	flag.Parse()

	// Set up the logrus logger:
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(logrus.DebugLevel)

	// Attempt to load config from file:
	conf.LoadFromFile()

	// Dump the config to the log:
	dumpConfig()
}

func dumpConfig() {
	// Dump the config:
	log.WithFields(logrus.Fields{"addresses": conf.MachineAddresses}).Debug("Config")
	log.WithFields(logrus.Fields{"concurrency": conf.Concurrency}).Debug("Config")
	log.WithFields(logrus.Fields{"configfile": conf.ConfFile}).Debug("Config")
	log.WithFields(logrus.Fields{"debug": conf.Debug}).Debug("Config")
	log.WithFields(logrus.Fields{"iterations": conf.Iterations}).Debug("Config")
	log.WithFields(logrus.Fields{"keepalive": conf.APIKeepAlive}).Debug("Config")
	log.WithFields(logrus.Fields{"playlist": conf.PlaylistFile}).Debug("Config")
	log.WithFields(logrus.Fields{"sleep": conf.SleepBetweenTests}).Debug("Config")
	log.WithFields(logrus.Fields{"testopenplaylist": conf.TestOpenPlaylist}).Debug("Config")
	log.WithFields(logrus.Fields{"testopenpopulateplaylist": conf.TestOpenPopulatePlaylist}).Debug("Config")
	log.WithFields(logrus.Fields{"testtriggerplaylist": conf.TestTriggerPlaylist}).Debug("Config")
	log.WithFields(logrus.Fields{"timeout": conf.APITimeout}).Debug("Config")
}

func main() {

	// Prepare a channel of results (to feed the resultprocessor):
	log.Info("Preparing the results channel")
	resultsChannel := make(chan types.Result)

	// Start the result processor in the background:
	go resultprocessor.Run(conf, resultsChannel, waitGroup)

	// Start enough testers to meet the concurrency parameter:
	for machineNumber := 0; machineNumber < conf.Concurrency; machineNumber++ {

		// Make a new tester:
		testers[machineNumber] = testaburger.Tester{
			Config:         conf,
			MachineNumber:  machineNumber,
			ResultsChannel: resultsChannel,
			WaitGroup:      waitGroup,
		}

		// Run it in the background:
		waitGroup.Add(1)
		go testers[machineNumber].Run()
	}

	// Make sure we wait for everything to complete before bailing out:
	waitGroup.Wait()

	// Report the results:
	resultprocessor.ReportSummary()
}
