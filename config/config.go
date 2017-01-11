package config

import (
	ioutil "io/ioutil"
	time "time"

	logrus "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var (
	log = logrus.WithFields(logrus.Fields{"logger": "Config"})
)

type Config struct {
	APIKeepAlive             time.Duration `yaml:"keepalive"`
	APITimeout               time.Duration `yaml:"timeout"`
	Concurrency              int           `yaml:"concurrency"`
	ConfFile                 string
	Debug                    bool          `yaml:"debug"`
	Iterations               int           `yaml:"iterations"`
	MachineAddresses         []string      `yaml:"addresses"`
	PlaylistFile             string        `yaml:"playlist"`
	SleepBetweenTests        time.Duration `yaml:"sleep"`
	TestOpenPlaylist         bool          `yaml:"testopenplaylist"`
	TestOpenPopulatePlaylist bool          `yaml:"testopenpopulateplaylist"`
	TestTriggerPlaylist      bool          `yaml:"testtriggerplaylist"`
}

func (myConfig *Config) LoadFromFile() {

	// Try to load the config file:
	yamlFile, err := ioutil.ReadFile(myConfig.ConfFile)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "filename": myConfig.ConfFile}).Debug("Couldn't load config-file")
		return
	}

	// Try to unmarshal the YAML into our config object:
	err = yaml.Unmarshal(yamlFile, myConfig)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "filename": myConfig.ConfFile}).Error("Couldn't parse config-file")
	}
}
