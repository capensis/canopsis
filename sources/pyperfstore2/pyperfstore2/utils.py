#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------
import logging
logger = logging.getLogger('utils')
logger.setLevel(logging.DEBUG)

import zlib
import time

import msgpack
packer = None
unpacker = None

import calendar
from datetime import datetime, timedelta
from dateutil.relativedelta import *

T_SECOND = 'second'
T_MINUTE = 'minute'
T_HOUR = 'hour'
T_DAY = 'day'
T_WEEK = 'week'
T_MONTH = 'month'
T_YEAR = 'year'

relativeDeltas = {
	T_MINUTE: relativedelta(minutes=1),
	T_HOUR: relativedelta(hours=1),
	T_DAY: relativedelta(days=1),
	T_WEEK: relativedelta(weeks=1),
	T_MONTH: relativedelta(months=1),
	T_YEAR: relativedelta(years=1)
}

S = 1
MN = 60
HR = 3600
D = 86400
W = 604800
M = 2629800
Y = 31557600

periodtypeByInterval = {
	S: T_SECOND,
	MN: T_MINUTE,
	HR: T_HOUR,
	D: T_DAY,
	W: T_WEEK,
	M: T_MONTH,
	Y: T_YEAR
}

intervalToRelativeDelta = {
	MN:relativedelta(minutes=+1),
	MN * 5:relativedelta(minutes=+5),
	MN * 15:relativedelta(minutes=+15),
	MN * 30:relativedelta(minutes=+30),
	HR:relativedelta(hours=+1),
	D:relativedelta(days=+1),
	W:relativedelta(weeks=+1),
	M:relativedelta(months=+1),
	Y:relativedelta(years=+1)
}

#### Utils fn
def datetimeToTimestamp(_date):
	return time.mktime(_date.timetuple())
	return calendar.timegm(_date.timetuple())

def get_overlap(a, b):
	return max(0, min(a[1], b[1]) - max(a[0], b[0]))

def get_first_point(points):
	if len(points):
		return points[0]
	else:
		return None

def get_last_point(points):
	if len(points):
		return points[len(points)-1]
	else:
		return None

def get_first_value(points):
	point = get_first_point(points)
	if point:
		return point[1]
	else:
		return None

def get_last_value(points):
	point = get_last_point(points)
	if point:
		return point[1]
	else:
		return None
		
def delta(points):
	if len(points) == 1:
		return points[0][1]
		
	vfirst = get_first_value(points)
	vlast = get_last_value(points)
	return vlast - vfirst

def median(vlist):
    values = sorted(vlist)
    count = len(values)

    if count % 2 == 1:
        return values[(count+1)/2-1]
    else:
        lower = values[count/2-1]
        upper = values[count/2]

    return (float(lower + upper)) / 2

def get_timestamp_interval(points):
	timestamp = 0
	timestamps=[]
	for point in points:
		timestamps.append(point[0] - timestamp)
		timestamp = point[0]

	if len(timestamps) > 1:
		del timestamps[0]

	return int(median(timestamps))

def get_timestamps(points):
	return [x[0] for x in points]

def get_values(points):
	return [x[1] for x in points]

def mean(vlist):
	if len(vlist):
		return round(float( sum(vlist) / float(len(vlist))), 3)
	else:
		return 0.0

def vmean(vlist):
	vlist = get_values(vlist)
	return mean(vlist)

def vmin(vlist):
	vlist = get_values(vlist)
	return min(vlist)

def vmax(vlist):
	vlist = get_values(vlist)
	return max(vlist)

def vsum(vlist):
	vlist = get_values(vlist)
	return sum(vlist)

def derivs(vlist):
	return [vlist[i] - vlist[i - 1] for i in range(1, len(vlist) - 2)]

def parse_dst(points, dtype, first_point=[]):
	logger.debug("Parse Data Source Type %s on %s points" % (dtype, len(points)))

	dtype = dtype.upper()
		
	if dtype == "DERIVE" or dtype == "COUNTER" or dtype == "ABSOLUTE":
		if points:
			rpoints = []
			values = get_values(points)
			i=0
			last_value=0
			counter = 0
			
			logger.debug('There is %s values' % len(values))
			
			for point in points:
				
				value = point[1]
				timestamp = point[0]
				
				previous_timestamp = None
				previous_value = None
				
				## Get previous value and timestamp
				if i != 0:
					previous_value 		= points[i-1][1]
					previous_timestamp	= points[i-1][0]
				elif i == 0 and first_point:
					previous_value		= first_point[1]
					previous_timestamp	= first_point[0]
				
				
				## Calcul Value
				if dtype != "COUNTER":
					if previous_value:
						if value > previous_value:
							value -= previous_value
						else:
							value = 0
				
				## Derive
				if previous_timestamp and dtype == "DERIVE":	
					interval = abs(timestamp - previous_timestamp)
					if interval:
						value = round(float(value) / interval, 3)
				
				## Abs
				if dtype == "ABSOLUTE":
					value = abs(value)
					
				## COUNTER
				if dtype == "COUNTER":
					value = value + counter
					counter = value

				## if new dca start, value = 0 and no first_point: wait second point ...
				if dtype == "DERIVE" and i == 0 and not first_point:
					## Drop this point
					pass
				else:
					rpoints.append([timestamp, value])
					
				i += 1
				
			return rpoints
	
	return points

def _roundtime(utcdate, periodtime=1, periodtype=T_HOUR, timezone=time.timezone):
	"""
	Calculate roudtime relative to an UTC date, a period time/type and a timezone.
	"""
	result = utcdate

	relativeinterval = intervalToRelativeDelta.get(periodtype, None)

	if relativeinterval != None:
		# assume result does not contain seconds and microseconds in this case
		result = result.replace(microsecond=0)

		if periodtype == T_SECOND:
			result = result.replace(second= 0 if periodtime > S else int(int(result.second / periodtime) * periodtime))
		elif periodtype == T_MINUTE:
			result = result.replace(second=0)
			result = result.replace(minute = 0 if periodtime > MN else int(int(result.minute / periodtime) * periodtime))
		elif periodtype == T_HOUR:
			result = result.replace(second=0, minute=0)
			result = result.replace(hour= 0 if periodtime > HR else int(int(result.hour / periodtime) * periodtime))
		elif periodtype == T_DAY:
			result = result.replace(second=0, minute=0, hour=0)
			#monthday = calendar.monthday(result.year, result.month)[1]
			#result = result.replace(day=1 if periodtime > monthday else (result.day / periodtime) * periodtime)
		elif periodtype == T_WEEK:
			result = result.replace(second=0, minute=0, hour=0)
			weeks = calendar.monthcalendar(result.year, result.month)
			for index in xrange(len(weeks)):
				week = weeks[index]
				if result.day in week:
					result = result.replace(day=week[0] if week[0] != 0 else 1)
					break			
		elif periodtype == T_MONTH:
			result = result.replace(second=0, minute=0, hour=0, day=1)			
		elif periodtype == T_YEAR:
			result = result.replace(second=0, minute=0, hour=0, day=1, month=1)

	td = timedelta(seconds=timezone)
	result += td

	return result

def roundTime(date, interval, timezone=time.timezone):
	"""
	Calculate roudtime relative to an UTC date, an interval.
	"""
	result = date

	relativeinterval = intervalToRelativeDelta.get(interval, None)

	if relativeinterval:

		dt = timedelta(seconds=timezone)
		result -= dt

		# assume result does not contain seconds and microseconds in this case
		result = result.replace(second=0, microsecond=0)

		if interval < HR: # in minutes
			minutes = (result.minute * 60 / interval) * interval / 60
			result = result.replace(minute=minutes)

		else:
			result = result.replace(minute=0)

			if interval >= D: # >= 1 day
				result = result.replace(hour=0)

			if interval >= W: # >= 1 week
				weeks = calendar.monthcalendar(result.year, result.month)
				for week in weeks:
					if result.day in week:
						result = result.replace(day=week[0] if week[0]!=0 else 1)
						break

			if interval >= M: # >= 1 month
				result = result.replace(day=1)

			if interval >= Y: # >= 1 year
				result = result.replace(month=1)

		result += dt

	return result

def _getTimeSteps(start, stop, periodtime, periodtype, roundtime, timezone=time.timezone):
	logger.debug('getTimeSteps:')
	timeSteps = []
	
	logger.debug('   + Interval: %s' % interval)

	start_datetime 	= datetime.utcfromtimestamp(start)
	stop_datetime 	= datetime.utcfromtimestamp(stop)

	if roundtime:
		stop_datetime = roundTime(stop_datetime, interval, timezone)

	if periodtype != None:

		relativeinterval = relativeDeltas[periodtype] * periodtime	

		if relativeinterval != None:
			date = stop_datetime
			start_datetime_minus_relativeinterval = start_datetime - relativeinterval
			while date > start_datetime_minus_relativeinterval:
				timeSteps.append(datetimeToTimestamp(date))
				date -= relativeinterval

	else:
		logger.debug('   + Use interval')
		timeSteps = range(stop, start-periodtime, -periodtime)
	
	timeSteps.reverse()
	
	logger.debug('   + timeSteps: ', timeSteps)

	return timeSteps

def getTimeSteps(start, stop, interval, roundtime=True, timezone=time.timezone):
	logger.debug('getTimeSteps:')
	timeSteps = []
	
	logger.debug('   + Interval: %s' % interval)

	start_datetime 	= datetime.utcfromtimestamp(start)
	stop_datetime 	= datetime.utcfromtimestamp(stop)

	if roundtime:
		stop_datetime = roundTime(stop_datetime, interval, timezone)

	relativeinterval = intervalToRelativeDelta.get(interval, None)

	if relativeinterval:
		date = stop_datetime
		start_datetime_minus_relativeinterval = start_datetime - relativeinterval

		while date > start_datetime_minus_relativeinterval:			
			ts = calendar.timegm(date.timetuple())			
			timeSteps.append(ts)
			date -= relativeinterval
	else:
		logger.debug('   + Use interval')
		timeSteps = range(stop, start-interval, -interval)
	
	timeSteps.reverse()
	
	logger.debug('   + timeSteps: ', timeSteps)

	return timeSteps

def aggregate(points, start=None, stop=None, max_points=None, interval=None, atype='MEAN', agfn=None, mode=None, fill=False, roundtime = True, timezone=time.timezone):

	if not atype:
		return points

	if not mode:
		mode = 'by_point'
	elif mode != 'by_point':
		mode = 'by_interval'
	
	if not max_points:
		max_points=1450
		
	if interval:
		interval = int(interval)
		mode = 'by_interval'
				
	if max_points != None:
		 max_points = int(max_points)

	atype = atype.upper()
	
	logger.debug("Aggregate %s points (max: %s, interval: %s, method: %s, mode: %s)" % (len(points), max_points, interval, atype, mode))

	if not agfn:
		if atype == 'MEAN':
			agfn = vmean
		elif atype == 'FIRST':
			agfn = get_first_value
		elif atype == 'LAST':
			agfn = get_last_value
		elif atype == 'MIN':
			agfn = vmin
		elif atype == 'MAX':
			agfn = vmax
		elif atype == 'DELTA':
			agfn = delta
		elif atype == 'SUM':
			agfn = vsum
		else:
			agfn = vmean

	logger.debug(" + Interval: %s" % interval)
	#logger.debug(" + Points: %s" % points)

	rpoints=[]
	
	if mode == 'by_point':
		if len(points) < max_points:
			logger.debug(" + Useless (%s < %s)" % (len(points), max_points))
			return points
		
		interval = int(round(len(points) / float(max_points)))
		logger.debug(" + point interval: %s" % interval)
		
		for x in range(0, len(points), interval):
			sample = points[x:x+interval]
			value = agfn(sample)
			timestamp = sample[len(sample)-1][0]
			rpoints.append([timestamp, value])
		
	elif mode == 'by_interval':
		
		if not start:
			start = points[0][0]

		if not stop:
			stop = points[len(points)-1][0]

		if len(points) == 1:
			return [ [start, points[0][1]] ]
		
		logger.debug(' + Start: %s' %  datetime.utcfromtimestamp(start))
		logger.debug(' + Stop:  %s' %  datetime.utcfromtimestamp(stop))

		timeSteps = getTimeSteps(start, stop, interval, roundtime, timezone)

		#initialize variables for loop
		prev_point = None
		i=0
		points_to_aggregate = []
		last_point = None

		for index in xrange(1, len(timeSteps)):

			timestamp = timeSteps[index]
			
			previous_timestamp = timeSteps[index-1]
			
			logger.debug("   + Interval %s -> %s" % (previous_timestamp, timestamp))

			while i < len(points) and points[i][0] < timestamp:

				points_to_aggregate.append(points[i])

				i+=1

			if atype == 'DELTA' and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_point = get_aggregation_point(points_to_aggregate, agfn, previous_timestamp, fill)

			rpoints.append(aggregation_point)

			if points_to_aggregate:
				last_point = points_to_aggregate[-1]

			points_to_aggregate = []

		if i < len(points):

			points_to_aggregate = points[i:]

			if atype == 'DELTA' and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_point = get_aggregation_point(points_to_aggregate, agfn, timeSteps[-1], fill)

			rpoints.append(aggregation_point)

	logger.debug(" + Nb points: %s" % len(rpoints))

	return rpoints

def get_aggregation_point(points_to_aggregate, fn, timestamp, fill):
	if points_to_aggregate:
		
		logger.debug("     + %s points to aggregate" % (len(points_to_aggregate)))

		agvalue = round(fn(points_to_aggregate), 2)

		result = [timestamp, agvalue]

	else:
		logger.debug("       + No points")

		result = [timestamp, 0 if fill else None]

	logger.debug("   + Point : %s " % result)

	return result

def compress(points):
	logger.debug("Compress timeserie")
	
	# Create packer
	global packer
	if not packer:
		packer = msgpack.Packer()
	
	# Remplace timestamp by interval
	logger.debug(" + Remplace Timestamp by Interval and compress it")
	i = 0
	fts = points[0][0]
	offset = points[0][0]
	previous_interval = None

	data = []
	
	logger.debug(" + FTS: %s" % fts)

	for point in points:
		timestamp = point[0]
		value = point[1]
		
		# If int, dont store float
		if value == int(value):
			value = int(value)
	
		if i == 0:
			# first point
			interval = timestamp - offset
			data.append(value)
		else:
			# Others
			interval = timestamp - offset
			if interval == previous_interval:
				data.append(value)
			else:
				previous_interval = interval
				data.append([interval, value])

		#logger.debug("    + %s: %s: %s" % (i, point, data[i]))
		
		offset = timestamp
		i += 1
	
	data = (fts, data)
	# Pack and compress points
	
	points = zlib.compress(packer.pack(data), 9)

	return points

def uncompress(data):
	logger.debug("Uncompress timeserie")
	
	if not data:
		raise ValueError("Invalid data type (%s)" % type(data))

	# Create unpacker
	global unpacker
	if not unpacker:
		unpacker = msgpack.Unpacker(use_list=True)
	
	unpacker.feed(str(zlib.decompress(data)))
	data = unpacker.unpack()
	
	fts = data[0]
	points = data[1]
	
	logger.debug(" + Type of point: %s" % type(points))
	
	if type(points).__name__ != 'list':
		raise ValueError("Invalid type (%s)" % type(points))
	
	rpoints = []

	#first point
	rpoints.append([fts, points[0]])
	logger.debug("   + First point: %s" % (rpoints[0]))

	#second point
	offset = points[1][0]
	timestamp = fts + offset
	rpoints.append([timestamp, points[1][1]])
	
	logger.debug("   + Second point: %s" % (rpoints[1]))
	
	logger.debug(" + Offset: %s", offset)

	#others
	for i in range(2, len(points)):
		point = points[i]
		
		if isinstance(point ,list) or isinstance(point ,tuple):
			offset = point[0]
			#logger.debug(" + Offset: %s", offset)
			timestamp += offset
			rpoints.append([ timestamp, point[1] ])
		else:
			timestamp += offset
			rpoints.append([ timestamp, point ])
			
		#logger.debug("  %i -> %s" % (i, rpoints[i]))
	
	return rpoints

### aggregation serie function
def consolidation(series, fn, interval=None):
	
	# Todo calcul interval
	if not interval:
		interval = 300

	# Find start and stop ts
	start = None
	stop = None
	for serie in series:
		timestamps = [point[0] for point in serie]
		smin = min(timestamps)
		smax = max(timestamps)
		if start == None or smin < start:
			start = smin
		if stop == None or smax > stop:
			stop = smax

	# Align timestamps
	nseries = []
	for serie in series:
		ts = start
		index = 0
		nserie = []
		last_value = None

		while ts <= stop:
			while index < len(serie) and serie[index][0] <= ts:
				last_value = serie[index][1]
				index += 1

			nserie.append( (ts, last_value) )
			ts+=interval

		nseries.append(nserie)

	# Do operations
	result = []
	ts = start
	i = 0
	while ts <= stop:
		points = []
		for serie in nseries:
			if serie[i][1]:
				points.append(serie[i][1])

		value = fn(points)
		result.append((ts, value))

		i  += 1
		ts += interval

	return result
