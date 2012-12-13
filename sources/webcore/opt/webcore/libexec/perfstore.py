#!/usr/bin/env python
# --------------------------------
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

import sys, os, logging, json, time
from datetime import datetime

import bottle
from bottle import route, get, post, put, delete, request, HTTPError, response

#import protection function
from libexec.auth import check_auth, get_account

# Modules
from ctools import parse_perfdata, clean_mfilter

import pyperfstore2
from pyperfstore2.utils import aggregate_series , mean
manager = None

import ConfigParser
config = ConfigParser.RawConfigParser()
config.read(os.path.expanduser('~/etc/cstorage.conf'))

logger = logging.getLogger("perfstore")

config.read(os.path.expanduser('~/etc/webserver.conf'))
pyperfstore_aggregate			= True
pyperfstore_aggregate_maxpoints	= 500
pyperfstore_aggregate_method	= "MAX"
try:
	pyperfstore_aggregate			= config.getboolean('pyperfstore', "aggregate")
	pyperfstore_aggregate_maxpoints	= config.getint('pyperfstore', "aggregate_maxpoints")
	pyperfstore_aggregate_method	= config.get('pyperfstore', "aggregate_method")
except:
	pass
	
logger.debug(" + pyperfstore_aggregate: %s" % pyperfstore_aggregate)
logger.debug(" + pyperfstore_aggregate_maxpoints: %s" % pyperfstore_aggregate_maxpoints)
logger.debug(" + pyperfstore_aggregate_method: %s" % pyperfstore_aggregate_method)

def load():
	global manager
	manager = pyperfstore2.manager(logging_level=logging.INFO)

def unload():
	global manager
	del manager

group_managing_access = ['group.CPS_perfdata_admin']
#########################################################################

#### POST@
@post('/perfstore/values')
@post('/perfstore/values/:start/:stop')
def perfstore_nodes_get_values(start=None, stop=None, interval=None):

	metas = request.params.get('nodes', default=None)
	
	aggregate_method	= request.params.get('aggregate_method',	default=None)
	aggregate_interval	= request.params.get('aggregate_interval', default=None)
	aggregate_max_points= request.params.get('aggregate_max_points', default=None)
	consolidation = request.params.get('consolidation', default=None)
	output = []
	
	if not metas:
		logger.warning("Invalid arguments")
		return HTTPError(404, "Invalid arguments")

	metas = json.loads(metas)
	
	logger.debug("POST:")
	logger.debug(" + metas: %s" % metas)
	logger.debug(" + aggregate_method: %s" % aggregate_method)
	logger.debug(" + aggregate_interval: %s" % aggregate_interval)
	logger.debug(" + aggregate_max_points: %s" % aggregate_max_points)

	output = []

	for meta in metas:
		_id = meta.get('id', None)
		if _id:
			output += perfstore_get_values(	_id=meta['id'],
											start=start,
											stop=stop,
											aggregate_method=aggregate_method,
											aggregate_interval=aggregate_interval,
											aggregate_max_points=aggregate_max_points)
	if consolidation and len(output) != 0:
		##select right function
		if consolidation == 'mean':
			fn = mean
		elif consolidation == 'min':
			fn = lambda x: min(x)
		elif consolidation == 'max' :
			fn = lambda x: max(x)
		elif consolidation == 'sum':
			fn = lambda x: sum(x)
		elif consolidation == 'delta':
			fn = lambda x: x[0] - x[-1]

		series = []
		for serie in output:
			series.append(serie["values"])
		output = [{
			'node': output[0]['node'],
			'metric': consolidation,
			'bunit': None,
			'type': 'GAUGE',
			'values': aggregate_series(series, fn, 60* 1000)
		}]

	output = {'total': len(output), 'success': True, 'data': output}
	return output
	
#### GET@
@get('/perfstore')
@get('/perfstore/get_all_metrics')
def perstore_get_all_metrics():
	logger.debug("perstore_get_all_metrics:")
	
	limit		= int(request.params.get('limit', default=20))
	start		= int(request.params.get('start', default=0))
	search		= request.params.get('search', default=None)
	filter		= request.params.get('filter', default=None)
	sort		= request.params.get('sort', default=None)
	show_internals	= request.params.get('show_internals', default=False)
	
	if filter:
		try:
			filter = json.loads(filter)
		except Exception, err:
			logger.error("Impossible to decode filter: %s: %s" % (filter, err))
			filter = None
	
	
	if show_internals == "true":
		show_internals = True
	else:
		show_internals = False
		
	msort = []
	if sort:
		sort = json.loads(sort)
		for item in sort:
			direction = 1
			if str(item['direction']) == "DESC":
				direction = -1
			msort.append((str(item['property']), direction))
	else:
		msort.append(('co', 1))
	
	logger.debug(" + limit:   %s" % limit)
	logger.debug(" + start:   %s" % start)
	logger.debug(" + search:  %s" % search)
	logger.debug(" + sort: "+str(sort))
	logger.debug(" + msort: "+str(msort))
	logger.debug(" + filter:  %s" % filter)
	logger.debug(" + show_internals:  %s" % show_internals)
	
	mfilter = None
	
	if isinstance(filter, list):
		if len(filter) > 0:
			mfilter = filter[0]
		else:
			logger.error(" + Invalid filter format")
			
	elif isinstance(filter, dict):
		mfilter = filter
	
	if search:
		# Todo: Tweak this ...
		fields = ['co', 're', 'me']
		mor = []
		search = search.split(' ')
		if len(search) == 1:
			for field in fields:
				mor.append({field: {'$regex': '.*%s.*' % search[0], '$options': 'i'}})
				
			mfilter = {'$or': mor}
		else:
			mfilter = {'$and': []}
			for word in search:
				mor = []
				for field in fields:
					mor.append({field: {'$regex': '.*%s.*' % word, '$options': 'i'}})	
				mfilter['$and'].append({'$or': mor})
	
	internal_metrics = [ 'cps_state', 'cps_statechange', 'cps_statechange_nok', 'cps_statechange_0', 'cps_statechange_1', 'cps_statechange_2', 'cps_statechange_3', 'cps_evt_per_sec', 'cps_sec_per_evt', 'cps_queue_size', 'cps_sel_state_0', 'cps_sel_state_1', 'cps_sel_state_2', 'cps_sel_state_3'  ]
	
	if not show_internals:
		if mfilter:
			mfilter = {'$and': [mfilter, {'me': {'$nin':internal_metrics  }}]}
		else:
			mfilter = {'me': {'$nin': internal_metrics  }}
		
	logger.debug(" + mfilter:  %s" % mfilter)
	
	mfilter = clean_mfilter(mfilter)

	data  = manager.find(limit=0, skip=0, mfilter=mfilter, data=False, sort=msort)
	total = data.count()
	data  = [meta for meta in data.skip(start).limit(limit)]
	
	return {'success': True, 'data' : data, 'total' : total}


### manipulating meta
@delete('/perfstore',checkAuthPlugin={'authorized_grp':group_managing_access})
@delete('/perfstore/:_id',checkAuthPlugin={'authorized_grp':group_managing_access})
def remove_meta(_id=None):
	if not _id:
		_id =  json.loads(request.body.readline())
	if not _id:
		return HTTPError(400, "No ids provided, bad request")
		
	if not isinstance(_id,list):
		_id = [_id]
	
	logger.debug('delete %s: ' % str(_id))
	
	for item in _id:
		if isinstance(item,dict):
			manager.remove(_id=item['_id'], purge=False)
		else:
			manager.remove(_id=item, purge=False)
			
@put('/perfstore',checkAuthPlugin={'authorized_grp':group_managing_access})
def update_meta(_id=None):
	data = json.loads(request.body.readline())
	
	if not isinstance(data,list):
		data = [data]
	
	for item in data:
		try:
			if not _id:
				_id = item['_id']
			if '_id' in item:
				del item['_id']
			manager.store.update(_id=_id, mset=item)
		except Exception, err:
			logger.warning('Error while updating meta_id: %s' % err)
			return HTTPError(500, "Error while updating meta_id: %s" % err)
			
"""@post('/perfstore/clean_all')
def clean_all():
	try:
		manager.cleanAll()
	except Exception,err:
		return HTTPError(500, "Error while cleaning all metrics: %s" % err)
"""
"""
@post('/perfstore/clean')
@post('/perfstore/clean/:_id')
def clean(_id=None):
	data = request.body.readline()
	if data:
		try:
			data = json.loads(data)
			logger.debug(data)
		except:
			logger.debug('No json decodable')
	
	if not isinstance(data,list):
		data = [data]
		
	if _id:
		logger.debug('Launch clean for: %s' % _id)
		manager.clean(_id=_id)
	else:
		for item in data:
			if isinstance(item, str) or isinstance(item, unicode):
				logger.debug('Launch clean for: %s' % item)
				manager.clean(_id=item)
			if isinstance(item, dict):
				logger.debug('Launch clean for: %s' % item['_id'])
				manager.clean(_id=item['_id'])
"""

#### POST@
@route('/perfstore/perftop')
def perfstore_perftop():
	data = []
	
	limit					= int(request.params.get('limit', default=10))
	sort					= int(request.params.get('sort', default=1))
	mfilter					= request.params.get('mfilter', default={})
	time_window				= int(request.params.get('time_window', default=86400))
	threshold				= request.params.get('threshold', default=None)
	threshold_direction 	= int(request.params.get('threshold_direction', default=-1))
	expand 					= request.params.get('expand', default=False)
	percent					= request.params.get('percent', default=False)
	threshold_on_pct		= request.params.get('threshold_on_pct', default=False)

	if percent == 'true':
		percent = True
	elif percent == 'false':
		percent = False

	if threshold_on_pct == 'true':
		threshold_on_pct = True
	elif threshold_on_pct == 'false':
		threshold_on_pct = False
	
	sort_on_percent = False
	if percent == True:
		sort_on_percent = True

	if mfilter:
		try:
			mfilter = json.loads(mfilter)
		except Exception, err:
			logger.error("Impossible to decode mfilter: %s: %s" % (mfilter, err))
			mfilter = None

	if threshold:
		threshold = float(threshold)

	if expand == 'true':
		expand = True
	else:
		expand = False

	logger.debug("PerfTop:")
	logger.debug(" + mfilter:     %s" % mfilter)
	logger.debug(" + limit:       %s" % limit)
	logger.debug(" + threshold:   %s" % threshold)
	logger.debug(" + threshold_direction:   %s" % threshold_direction)
	logger.debug(" + time_window: %s" % time_window)
	logger.debug(" + sort:        %s" % sort)
	logger.debug(" + expand:       %s" % expand)
	logger.debug(" + percent:       %s" % percent)
	logger.debug(" + threshold_on_pct:       %s" % threshold_on_pct)

	mfilter =  clean_mfilter(mfilter)
	
	mtype = manager.store.find(mfilter=mfilter, limit=1, mfields=['t'])
	
	def check_threshold(value):
		if threshold:
			if threshold_direction == -1 and value >= threshold:
				return True
			elif threshold_direction == 1 and value <= threshold:
				return True
			else:
				return False
		else:
			return True

	if mtype:
		mtype = mtype.get('t', 'GAUGE')

		logger.debug(" + mtype:    %s" % mtype)
		
		if mtype != 'COUNTER' and not expand:
			metrics = manager.store.find(mfilter=mfilter, mfields=['_id', 'co', 're', 'me', 'lv', 'u', 'ma', 'lts'], sort=[('lv', sort)], limit=limit)
			
			if isinstance(metrics, dict):
				metrics = [metrics]

			for metric in metrics:
				if (percent or threshold_on_pct) and 'ma' in metric and 'lv' in metric:
					metric['pct'] = round(((metric['lv'] * 100)/ metric['ma']) * 100) / 100
				
				if threshold_on_pct:
					val = metric['pct']
				else:
					val = metric['lv']

				if check_threshold(val):
					data.append(metric)
		else:
			# Compute values
			tstop = int(time.time())
			tstart = tstop - time_window
			logger.debug(" + tstart:      %s" % tstart)
			logger.debug(" + tstop:       %s" % tstop)
			metric_limit = 0
			
			if expand:
				metric_limit = 1

			#clean mfilter
			mfilter =  clean_mfilter(mfilter)

			metrics =  manager.store.find(mfilter=mfilter, mfields=['_id', 'co', 're', 'me', 'lv', 'u', 'ma', 'lts'], limit=metric_limit)

			if isinstance(metrics, dict):
				metrics = [metrics]

			for metric in metrics:
				# Recheck type
				mtype = metric.get('t', 'GAUGE')

				if mtype != 'COUNTER' and not expand:
					logger.debug(" + Metric '%s' (%s) is not a COUNTER" % (metric['me'], metric['_id']))

					if (percent or threshold_on_pct) and 'ma' in metric and 'lv' in metric:
							metric['pct'] = round(((metric['lv'] * 100)/ metric['ma']) * 100) / 100

					if threshold_on_pct:
						val = metric['pct']
					else:
						val = metric['lv']

					if check_threshold(val):
						data.append(metric)
				else:
					points = manager.get_points(_id=metric['_id'], tstart=tstart, tstop=tstop)
					if expand:
						del metric['_id']
						if not len(points):
							metric['lv'] = 0
							data.append(metric)
						else:
							for point in points:
								if check_threshold(point[1]):
									nmetric = metric.copy()
									nmetric['lts'] = point[0]
									nmetric['lv'] = point[1]
									if (percent or threshold_on_pct) and 'ma' in nmetric and 'lv' in nmetric:
										nmetric['pct'] = round(((nmetric['lv'] * 100)/ nmetric['ma']) * 100) / 100
									data.append(nmetric)
					else:
						if len(points):
							metric['lv'] = points[len(points)-1][1]
						else:
							metric['lv'] = 0

						if (percent or threshold_on_pct) and 'ma' in metric and 'lv' in metric:
								metric['pct'] = round(((metric['lv'] * 100)/ metric['ma']) * 100) / 100

						if threshold_on_pct:
							val = metric['pct']
						else:
							val = metric['lv']

						if check_threshold(val):
							data.append(metric)
				
		reverse = True
		if sort == 1:
			reverse = False	

		if sort_on_percent:
			data = sorted(data, key=lambda k: k['pct'], reverse=reverse)[:limit]
		else:
			data = sorted(data, key=lambda k: k['lv'], reverse=reverse)[:limit]
	else:
		logger.debug("No records found")
		
	
	return {'success': True, 'data' : data, 'total' : len(data)}

########################################################################
# Functions
########################################################################

def perfstore_get_values(_id, start=None, stop=None, aggregate_method=None, aggregate_interval=None, aggregate_max_points=None):
	
	if start and not stop:
		stop = start
	
	if stop:
		stop = int(int(stop) / 1000)
	else:
		stop = int(time.time())
		
	if start:
		start = int(int(start) / 1000)
	else:
		start = stop - 86400

	if aggregate_interval:
		aggregate_interval = int(aggregate_interval)
	
	if not aggregate_max_points:
		aggregate_max_points = pyperfstore_aggregate_maxpoints
		
	if not aggregate_method:
		aggregate_method = pyperfstore_aggregate_method
	
	logger.debug("Perfstore get points:")
	logger.debug(" + meta _id:    %s" % _id)
	logger.debug(" + start:       %s (%s)" % (start, datetime.utcfromtimestamp(start)))
	logger.debug(" + stop:        %s (%s)" % (stop, datetime.utcfromtimestamp(stop)))
	logger.debug('Aggregate:')
	logger.debug(' + method :     %s' % aggregate_method)
	logger.debug(' + interval :   %s' % aggregate_interval)
	logger.debug(' + max_points : %s' % aggregate_max_points)
	
	output=[]
	meta = None
	
	if not _id:
		logger.error("Invalid _id '%s'" % _id)
		return output
	
	if (aggregate_interval):
		start -= start % aggregate_interval
		stop -= (stop % aggregate_interval)
		stop += aggregate_interval
		
		aggregate_max_points = int( round((stop - start) / aggregate_interval + 0.5) )
		logger.debug(" + start:       %s" % start)
		logger.debug(" + stop:        %s" % stop)
	
	try:
		points = []
		
		if start == stop:
			# Get only one point
			logger.debug("   + Get one point at %s: %s" % (stop, datetime.utcfromtimestamp(start)))
			(meta, point) = manager.get_point(	_id=_id,
												ts=start,
												return_meta=True)
			if point:
				points = [ point ]
				
			logger.debug('Point: %s' % points)
				
		else:
			
			(meta, points) = manager.get_points(	_id=_id,
												tstart=start,
												tstop=stop,
												return_meta=True)

			# For UI display
			if len(points) == 0 and meta['type'] == 'COUNTER':
				points = [(start, 0), (stop, 0)]
				
			points =  pyperfstore2.utils.aggregate(	points=points,
												max_points=aggregate_max_points,
												interval=aggregate_interval,
												atype=aggregate_method)
	except Exception, err:
		logger.error("Error when getting points: %s" % err)
	
	
	if(aggregate_interval):
		points = pyperfstore2.utils.fill_interval(points,start,stop,aggregate_interval)
	
	if points and meta:
		points = [[point[0] * 1000, point[1]] for point in points]

		output.append({'node': _id, 'metric': meta['me'], 'values': points, 'bunit': meta['unit'], 'min': meta['min'], 'max': meta['max'], 'thld_warn': meta['thd_warn'], 'thld_crit': meta['thd_crit'], 'type': meta['type']})
				
	return output
	
	
