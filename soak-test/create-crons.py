#!/usr/bin/env python
import sys

cron_user = 'root'

# Open the schedules CSV and split by '\n' (new-line):
with open(sys.argv[1]) as schedules_file:
    schedules = schedules_file.read().split('\n')

    # If we only got 1 line then assume that we split with the wrong character, try again with '\r' (carriage-return, the Windows default):
    if len(schedules) == 1:
		schedules_file.seek(0)
		schedules = schedules_file.read().split('\r')


# Iterate through it:
for schedule in schedules:

	try:
		# Split the schedule up (by comma):
		schedule_parts = schedule.split(',')

		# Split the time-component up:
		cron_hour =   schedule_parts[0].split(':')[0]
		cron_minute = schedule_parts[0].split(':')[1]
		cron_hostname = schedule_parts[6]
		cron_playlist = schedule_parts[9]
		cron_command = '/usr/local/bin/scheduled-test.sh %s %s' % (cron_hostname, cron_playlist)

		print('# %s' % schedule.strip('\n'))
		print('%s %s * * *\t%s\t%s\n' % (cron_minute, cron_hour, cron_user, cron_command))
	except Exception, e:
		pass
		# print('Exception: %s', e)
