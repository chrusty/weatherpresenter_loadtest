
WeatherPresenter load-testing tool
==================================

Parameters
----------
* __debug__: Run in DEBUG mode [_false_]
* __conffile__: Name of a config-file to load [_loadtest.yaml_]


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
