
MeteoGroup WeatherPresenter load-testing tool
=============================================

Parameters
----------
* __concurrency__: How many concurrent tests to run [_1_]
* __debug__: Run in DEBUG mode [_false_]
* __iterations__: Number of times to run each tests [_10_]
* __keepalive__: How often to send keepalive packets [_5s_]
* __addresses__: Comma-delimited list of WeatherPresenter web-remote API addresses to use [_http://localhost:34567_]
* __playlist__: Full path to a playlist to use
* __testloadplaylist__: Run the "load playlist" test (simply loads the playlist from disk) [_true_]
* __testpopulateplaylist__: Result [_false_]
* __timeout__: How long to wait for connections before timing out [_300s_]


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
