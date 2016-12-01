package config

import (
	time "time"
)

type Config struct {
	APIKeepAlive         time.Duration
	APITimeout           time.Duration
	Concurrency          int
	Iterations           int
	MachineAddresses     string
	PlaylistFile         string
	TestLoadPlaylist     bool
	TestPopulatePlaylist bool
}
