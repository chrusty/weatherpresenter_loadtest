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
	flag.IntVar(&conf.Iterations, "iterations", 0, "Number of times to run each tests (0 = forever)")
	flag.StringVar(&conf.MachineAddresses, "addresses", "http://playout-1:34567,http://playout-2:34567", "Comma-delimited list of WeatherPresenter machines to use")
	flag.StringVar(&conf.PlaylistFile, "playlist", `\\server\playlists\playlist.dlp`, "Full path to a playlist to use")
	flag.BoolVar(&conf.TestLoadPlaylist, "testloadplaylist", false, "Run the 'load playlist' test (simply loads the playlist from disk)")
	flag.BoolVar(&conf.TestPopulatePlaylist, "testpopulateplaylist", true, "Run the 'populate playlist' test (loads the playlist from disk, switches to 'Edit' mode, then closes the playlist)")
	flag.Parse()

	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)

	// Dump the config to the log:
	dumpConfig()
}

func dumpConfig() {
	// Dump the config:
	log.WithFields(logrus.Fields{"addresses": conf.MachineAddresses}).Debug("Config")
	log.WithFields(logrus.Fields{"concurrency": conf.Concurrency}).Debug("Config")
	log.WithFields(logrus.Fields{"iterations": conf.Iterations}).Debug("Config")
	log.WithFields(logrus.Fields{"keepalive": conf.APIKeepAlive}).Debug("Config")
	log.WithFields(logrus.Fields{"playlist": conf.PlaylistFile}).Debug("Config")
	log.WithFields(logrus.Fields{"testloadplaylist": conf.TestLoadPlaylist}).Debug("Config")
	log.WithFields(logrus.Fields{"testpopulateplaylist": conf.TestPopulatePlaylist}).Debug("Config")
	log.WithFields(logrus.Fields{"timeout": conf.APITimeout}).Debug("Config")
}

func main() {

	// Make sure we wait for everything to complete before bailing out:
	defer waitGroup.Wait()

	// Prepare a channel of results (to feed the resultprocessor):
	log.Info("Preparing the results channel")
	resultsChannel := make(chan types.Result)

	// Prepare to have background GoRoutines running:
	waitGroup.Add(1)

	// Start webhook server:
	go resultprocessor.Run(conf, resultsChannel, waitGroup)

	// Start enough testers to meet the concurrency parameter:
	for machineNumber := 0; machineNumber < conf.Concurrency; machineNumber++ {

		// Make a new tester:
		testers[machineNumber] = testaburger.Tester{
			Config:         conf,
			MachineNumber:  machineNumber,
			ResultsChannel: resultsChannel,
		}

		// Run it in the background:
		go testers[machineNumber].Run()
	}

}
