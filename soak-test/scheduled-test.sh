#!/usr/bin/env bash
LOG_FILE='/var/log/scheduled-tests.log'
TEST_ID=$$
WP_HOSTNAME=$1
WP_PLAYLISTINDEX=$2
DISABLEFILE='/tmp/scheduled-tests.disabled'
WEATHERPRESENTER_API_ADDRESS="http://${WP_HOSTNAME}:34567/weatherpresenter"
OPERATION_TIMEOUT=1800
KEEPALIVE_TIME=5

# Quit if we find a known file:
if [ -f $DISABLEFILE ]
then
	echo "`date`,[ERROR],,,,Tests are disabled (presence of ${DISABLEFILE}) - Exiting ..." >>$LOG_FILE
	exit
else
	echo "`date`,[INFO],,,,Scheduled test begins ..." >>$LOG_FILE
fi

# Make sure we've been given a parameter:
if [ $# -lt 1 ]
then
	echo "`date`,[ERROR],,,,Please provide a parameter (WP_HOSTNAME)! Exiting ..." >>$LOG_FILE
	exit
fi

# Figure out which playlist to load:
case $WP_PLAYLISTINDEX in
0)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_Production%5CPlaylists%5CTechnical+Tests%5CRegional+Birmingham.dpl'
  ;;
1)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_Production%5CPlaylists%5CTechnical+Tests%5CScotland+Evening.dpl'
  ;;
2)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_Production%5CPlaylists%5CTechnical+Tests%5CUK+1030+news.dpl'
  ;;
3)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_Production%5CPlaylists%5CTechnical+Tests%5CUK+1830+news.dpl'
  ;;
4)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_automatedTest%5CPlaylists%5CTechnical+Tests%5CHeavy+Show.dpl'
  ;;
*)
  PLAYLIST_FILEPATH='%5C%5CFGBW1APWTRB002%5Cmyshare%5CSoftware%5CDesignModel%5CBBC_automatedTest%5CPlaylists%5CTechnical+Tests%5CLight+Show.dpl'
  ;;
esac

# Record the time:
TIME_BEGIN=`date +%s`

# Close the current playlist, then sleep to allow the action to complete:
curl -s --max-time $OPERATION_TIMEOUT --keepalive-time $KEEPALIVE_TIME "${WEATHERPRESENTER_API_ADDRESS}/ClosePlaylist" >/dev/null
if [ "$RC" != 0 ]
then
	TIME_END=`date +%s`
	TIME_ELAPSED=$((TIME_END - TIME_BEGIN))
	echo "`date`,[ERROR],${WP_HOSTNAME},${TEST_ID},${TIME_ELAPSED}s,Couldn't close playlist!" >>$LOG_FILE
	exit
else
	sleep 10
fi

# Open the test playlist, then sleep to allow the action to complete:
curl -s --max-time $OPERATION_TIMEOUT --keepalive-time $KEEPALIVE_TIME "${WEATHERPRESENTER_API_ADDRESS}/OpenPlaylist?filepath=${PLAYLIST_FILEPATH}" >/dev/null
if [ "$RC" != 0 ]
then
	TIME_END=`date +%s`
	TIME_ELAPSED=$((TIME_END - TIME_BEGIN))
	echo "`date`,[ERROR],${WP_HOSTNAME},${TEST_ID},${TIME_ELAPSED}s,Couldn't open playlist!" >>$LOG_FILE
	exit
else
	sleep 15
fi

# Switch to "edit" mode (and record the time):
curl -s --max-time $OPERATION_TIMEOUT --keepalive-time $KEEPALIVE_TIME "${WEATHERPRESENTER_API_ADDRESS}/SetPresentationState?state=Editing" >/dev/null
RC=$?
TIME_END=`date +%s`
TIME_ELAPSED=$((TIME_END - TIME_BEGIN))
if [ "$RC" != 0 ]
then
	echo "`date`,[ERROR],${WP_HOSTNAME},${TEST_ID},${TIME_ELAPSED}s,Couldn't change to 'Editing' presentation mode!" >>$LOG_FILE
else
	echo "`date`,[SUCCESS],${WP_HOSTNAME},${TEST_ID},${TIME_ELAPSED}s,Finished" >>$LOG_FILE
fi
