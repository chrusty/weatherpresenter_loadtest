
WeatherPresenter load-testing tool
==================================

Parameters
----------
* __debug__: Run in DEBUG mode [_false_]
* __conffile__: Name of a config-file to load [_loadtest.yaml_]

Config
------
* __concurrency__: How many concurrent tests to run
* __debug__: Run in DEBUG mode
* __iterations__: Number of times to run each tests
* __conffile__: Name of a config-file to load
* __keepalive__: How often to send keepalive packets
* __addresses__: A list of hostnames / addresses of WeatherPresenter workstations to test against
* __playlist__: Full path to a playlist to use
* __sleep__: How long to sleep after running each test
* __testopenplaylist__: Run the "load playlist" test (simply loads the playlist from disk)
* __testopenpopulateplaylist__: Run the 'open & populate playlist' test (closes the playlist loads the playlist from disk, sleeps, switches to 'Edit' mode)
* __testtriggerplaylist__: Triggers REWIND then PLAY on the currently-loaded playlist, then sleeps
* __timeout__: How long to wait for connections before timing out

How to use
----------
* Make a config file
* Run the loadtest and direct stdout to a .csv file: `weatherpresenter_loadtest >results.csv`


WeatherPresenter Web-remote API
-------------------------------

### Playout control (from presentation or editing state)
* http://localhost:34567/weatherpresenter/trigger
* http://localhost:34567/weatherpresenter/forward
* http://localhost:34567/weatherpresenter/rewind
* http://localhost:34567/weatherpresenter/reset
* http://localhost:34567/weatherpresenter/back
* http://localhost:34567/weatherpresenter/refresh
* http://localhost:34567/weatherpresenter/corner
* http://localhost:34567/weatherpresenter/jumpTo?segment={segment}

### Playlist control
#### Open a playlist (filepath must be URL-encoded)
* http://localhost:34567/weatherpresenter/OpenPlaylist?filepath={filepath}

#### Close current playlist
* http://localhost:34567/weatherpresenter/ClosePlaylist

### Presentation-state
```
http://localhost:34567/weatherpresenter/SetPresentationState?state={None}
```

#### Options
* None
* Editing
* Recording 
