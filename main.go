package main

import (
	flag "flag"
	sync "sync"
	time "time"

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
	flag.DurationVar(&conf.APIKeepAlive, "keepalive", 5*time.Second, "How often to send keepalive packets")
	flag.DurationVar(&conf.APITimeout, "timeout", 300*time.Second, "How long to wait for connections before timing out")
	flag.IntVar(&conf.Concurrency, "concurrency", 1, "Number of concurrent tests to run")
	flag.BoolVar(&conf.Debug, "debug", false, "Run in DEBUG mode")
	flag.IntVar(&conf.Iterations, "iterations", 10, "Number of times to run each tests")
	flag.StringVar(&conf.MachineAddresses, "addresses", "http://localhost:34567", "Comma-delimited list of WeatherPresenter machines to use")
	flag.StringVar(&conf.PlaylistFile, "playlist", `\\server\playlists\playlist.dlp`, "Full path to a playlist to use")
	flag.DurationVar(&conf.SleepBetweenTests, "sleep", 2*time.Second, "How long to sleep after running each test")
	flag.BoolVar(&conf.TestOpenPlaylist, "testopenplaylist", false, "Run the 'open playlist' test (simply loads the playlist from disk)")
	flag.BoolVar(&conf.TestOpenPopulatePlaylist, "testopenpopulateplaylist", false, "Run the 'open & populate playlist' test (closes the playlist loads the playlist from disk, sleeps, switches to 'Edit' mode)")
	flag.BoolVar(&conf.TestTriggerPlaylist, "testtriggerplaylist", false, "Triggers REWIND then PLAY on the currently-loaded playlist, then sleeps")
	flag.Parse()

	// Set up the logrus logger:
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(logrus.DebugLevel)

	// Dump the config to the log:
	dumpConfig()
}

func dumpConfig() {
	// Dump the config:
	log.WithFields(logrus.Fields{"addresses": conf.MachineAddresses}).Debug("Config")
	log.WithFields(logrus.Fields{"concurrency": conf.Concurrency}).Debug("Config")
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
	resultprocessor.ProduceCSV()
}
