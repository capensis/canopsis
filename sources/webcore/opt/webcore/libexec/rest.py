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

import sys, os, logging, json

import bottle
from bottle import route, get, put, delete, request, HTTPError, post, response
from datetime import *
## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord
import base64
from ctools import clean_mfilter

#import protection function
from libexec.auth import get_account ,check_group_rights

logger = logging.getLogger("rest")

ctype_to_group_access = {
							'schedule' : 'group.CPS_schedule_admin',
							'curve' : 'group.CPS_curve_admin',
							'account' : 'group.CPS_account_admin',
							'group' : 'group.CPS_account_admin',
							'selector' : 'group.CPS_selector_admin',
							'derogation' : 'group.CPS_derogation_admin',
							'consolidation' : 'group.CPS_consolidation_admin'
						}

#########################################################################

#### GET Media
@get('/rest/media/:namespace/:_id')
def rest_get_media(namespace, _id):
	account = get_account()
	storage = get_storage(namespace=namespace)
	
	logger.debug("Get media '%s' from '%s':" % (_id, namespace))
	
	try:
		raw = storage.get(_id, account=account, namespace=namespace, mfields=["media_bin", "media_name", "media_type"], ignore_bin=False)
		
		media_type = raw.get('media_type', None)
		media_name = raw.get('media_name', None)
		media_bin = raw.get('media_bin', None)
		
	except Exception, err:
		logger.error(err)
		return HTTPError(404, err)
	
	if not media_type or not media_name or not media_bin:
		logger.error("Insufficient field in record")
		return HTTPError(404, "Insufficient field in record")
		
	logger.debug(" + media_type: %s" % media_type)	
	logger.debug(" + media_name: %s" % media_name)
	logger.debug(" + media_bin:  %s" % len(media_bin))
	
	response.headers['Content-Disposition'] = 'attachment; filename="%s"' % media_name
	response.headers['Content-Type'] = media_type
	
	return base64.b64decode(media_bin)

#### GET
@get('/rest/:namespace/:ctype/:_id')
@get('/rest/:namespace/:ctype')
@get('/rest/:namespace')
def rest_get(namespace, ctype=None, _id=None):
	#get the session (security)
	account = get_account()

	limit		= int(request.params.get('limit', default=20))
	page		= int(request.params.get('page', default=0))
	start		= int(request.params.get('start', default=0))
	groups		= request.params.get('groups', default=None)
	search		= request.params.get('search', default=None)
	filter		= request.params.get('filter', default=None)
	sort		= request.params.get('sort', default=None)
	query		= request.params.get('query', default=None)
	onlyWritable	= request.params.get('onlyWritable', default=False)
	noInternal	= request.params.get('noInternal', default=False)
	ids			= request.params.get('ids', default=[])
	
	get_id			= request.params.get('_id', default=None)
	
	fields = request.params.get('fields', default=None)
	
	if not _id and get_id:
		_id  = get_id
	
	if not isinstance(ids, list):
		try:
			ids = json.loads(ids)
		except Exception, err:
			logger.error("Impossible to decode ids: %s: %s" % (ids, err))

	if filter:
		try:
			filter = json.loads(filter)
		except Exception, err:
			logger.error("Impossible to decode filter: %s: %s" % (filter, err))
			filter = None
			

	msort = []
	if sort:
		#[{"property":"timestamp","direction":"DESC"}]
		sort = json.loads(sort)
		for item in sort:
			direction = 1
			if str(item['direction']) == "DESC":
				direction = -1

			msort.append((str(item['property']), direction))
		

	logger.debug("GET:")
	logger.debug(" + User: "+str(account.user))
	logger.debug(" + Group(s): "+str(account.groups))
	logger.debug(" + namespace: "+str(namespace))
	logger.debug(" + Ctype: "+str(ctype))
	logger.debug(" + _id: "+str(_id))
	logger.debug(" + ids: "+str(ids))
	logger.debug(" + Limit: "+str(limit))
	logger.debug(" + Page: "+str(page))
	logger.debug(" + Start: "+str(start))
	logger.debug(" + Groups: "+str(groups))
	logger.debug(" + onlyWritable: "+str(onlyWritable))
	logger.debug(" + Sort: "+str(sort))
	logger.debug(" + MSort: "+str(msort))
	logger.debug(" + Search: "+str(search))
	logger.debug(" + filter: "+str(filter))
	logger.debug(" + query: "+str(query))
	
	storage = get_storage(namespace=namespace)
	
	total = 0
	
	mfilter = {}
	if isinstance(filter, list):
		if len(filter) > 0:
			mfilter = filter[0]
		else:
			logger.error(" + Invalid filter format")
			
	elif isinstance(filter, dict):
		mfilter = filter
	
	records = []
	if ctype:
		if mfilter:
			mfilter['crecord_type'] = ctype
		else:
			mfilter = {'crecord_type': ctype}
			
	if query:
		if mfilter:
			mfilter['crecord_name'] = { '$regex' : '.*'+str(query)+'.*', '$options': 'i' }
		else:
			mfilter = {'crecord_name': { '$regex' : '.*'+str(query)+'.*', '$options': 'i' }}


	if _id:
		ids = _id.split(',')
		
	if ids:	
		try:
			records = storage.get(ids, account=account)
			if isinstance(records,crecord):
				records = [records]
				total = 1
			elif isinstance(records,list):
				total = len(records)
			else:
				total = 0
		except Exception, err:
			logger.info('Error: %s' % err)
			total = 0
		
		if total == 0:
			return HTTPError(404, str(ids) +" Not Found")
						
	else:
		if search:
			mfilter['_id'] = { '$regex' : '.*'+search+'.*', '$options': 'i' }
		
		logger.debug(" + mfilter: "+str(mfilter))

		#clean mfilter
		mfilter = clean_mfilter(mfilter)

		records =  storage.find(mfilter, sort=msort, limit=limit, offset=start, account=account)
		total =	storage.count(mfilter, account=account)

	output = []
	
	#----------------dump record and post filtering-------
	for record in records:
		if record:
			do_dump = True

			if onlyWritable:
				if not record.check_write(account=account):
					do_dump = False

			if noInternal:
				if 'internal' in record.data and record.data['internal']:
					do_dump = False

			if do_dump:
				data = record.dump(json=True)
				data['id'] = data['_id']
				if data.has_key('next_run_time'):
					data['next_run_time'] = str(data['next_run_time'])
				
				#Clean non wanted field
				if fields:
					fields_to_delete = []
					for item in data:
						if not item in fields:
							fields_to_delete.append(item)
					for field in fields_to_delete:
						del data[field]
							
					
				output.append(data)

	output={'total': total, 'success': True, 'data': output}
	
	return output
	
#### POST
@post('/rest/:namespace/:ctype/:_id')
@post('/rest/:namespace/:ctype')
def rest_post(namespace, ctype, _id=None):
	#get the session (security)
	account = get_account()
	storage = get_storage(namespace=namespace)
	
	#check rights on specific ctype (check ctype_to_group_access variable below)
	if ctype in ctype_to_group_access:
		if not check_group_rights(account,ctype_to_group_access[ctype]):
			return HTTPError(403, 'Insufficient rights')

	logger.debug("POST:")

	data = request.body.readline()
	if not data:
		return HTTPError(400, "No data received")

	logger.debug(" + data: %s" % data)
	logger.debug(" + data-type: %s" % type(data))
		
	if isinstance(data, str):
		try:
			data = json.loads(data)
		except Exception, err:
			logger.error("PUT: Impossible to parse data (%s)" % err)
			return HTTPError(404, "Impossible to parse data")

	data['crecord_type'] = ctype
	
	if not _id:
		try:
			_id = str(data['_id'])
		except:
			pass

		try:
			_id = str(data['id'])
		except:
			pass
	
	## Clean data
	try:
		del data['_id']
	except:
		pass

	try:
		del data['id']
	except:
		pass
	
	logger.debug(" + _id: "+str(_id))
	logger.debug(" + ctype: "+str(ctype))
	logger.debug(" + Data: "+str(data))
	
	## Set group
	if data.has_key('aaa_group'):
		group = data['aaa_group']
	else:
		group = account.group

	record = None
	if _id:
		try:
			record = storage.get(_id ,account=account)
			logger.debug('Update record %s' % _id)
		except:
			logger.debug('Create record %s' % _id)

	if record:
		for key in dict(data).keys():
			record.data[key] = data[key]
			
		# Update Name	
		try:
			record.name = data['crecord_name']
		except:
			pass
		
	else:
		raw_record = crecord(_id=_id, type=str(ctype)).dump()
		for key in dict(data).keys():
			raw_record[key] = data[key]

		record = crecord(raw_record=raw_record)
		record.chown(account.user)
		record.chgrp(group)
		#if ctype in ctype_to_group_access:
			#record.admin_group = ctype_to_group_access[ctype]
	
	logger.debug(' + Record: %s' % record.dump())
	try:
		storage.put(record, namespace=namespace, account=account)
		
	except Exception, err:
		logger.error('Impossible to put (%s)' % err)
		return HTTPError(403, "Access denied")
		
#### PUT
@put('/rest/:namespace/:ctype/:_id')
@put('/rest/:namespace/:ctype')
def rest_put(namespace, ctype, _id=None):
	#get the session (security)
	account = get_account()
	storage = get_storage(namespace=namespace)
	
	#check rights on specific ctype (check ctype_to_group_access variable below)
	if ctype in ctype_to_group_access:
		if not check_group_rights(account,ctype_to_group_access[ctype]):
			return HTTPError(403, 'Insufficient rights')

	logger.debug("PUT:")

	data = request.body.readline()
	if not data:
		return HTTPError(400, "No data received")

	logger.debug(" + data: %s" % data)
	logger.debug(" + data-type: %s" % type(data))
		
	if isinstance(data, str):
		try:
			data = json.loads(data)
		except Exception, err:
			logger.error("PUT: Impossible to parse data (%s)" % err)
			return HTTPError(404, "Impossible to parse data")
			
	if not _id:
		try:
			_id = str(data['_id'])
		except:
			pass

		try:
			_id = str(data['id'])
		except:
			pass
	
	## Clean data
	try:
		del data['_id']
	except:
		pass

	try:
		del data['id']
	except:
		pass
	
	logger.debug(" + _id: "+str(_id))
	logger.debug(" + ctype: "+str(ctype))
	logger.debug(" + Data: "+str(data))
	
	try:
		storage.update(_id, data, namespace=namespace, account=account)
		
	except Exception, err:
		logger.error('Impossible to put (%s)' % err)
		return HTTPError(403, "Access denied")


#### DELETE
@delete('/rest/:namespace/:ctype/:_id')
@delete('/rest/:namespace/:ctype')
def rest_delete(namespace, ctype, _id=None):
	account = get_account()
	storage = get_storage(namespace=namespace)
	
	logger.debug("DELETE:")
	
	data = request.body.readline()
	if data:
		try:
			data = json.loads(data)
		except:
			logger.warning('Invalid data in request payload')
			data = None	

	if data:
		logger.debug(" + Data: %s" % data)

		if isinstance(data, list):
			logger.debug(" + Attempt to remove %i item from db" % len(data))
			_id = []
				
			for item in data:
				if isinstance(item,str):
					_id.append(item)
					
				if isinstance(item,dict):
					item_id = item.get('_id', item.get('id', None))
					if item_id:
						_id.append(item_id)

		if isinstance(data, str):
			_id = data

		if isinstance(data, dict):
			_id = data.get('_id', data.get('id', None))

	if not _id:
		logger.error("DELETE: No '_id' field in header ...")
		return HTTPError(404, "No '_id' field in header ...")

	logger.debug(" + _id: %s" % _id)
	
	try:
		storage.remove(_id, account=account)
	except:
		return HTTPError(404, _id+" Not Found")



#decorateur, @get for http method, and then routes
@get('/event_log')
@get('/event_log/:start')
@get('/event_log/:start/:stop')
def get_timeline(start=None,stop=None):
	#get user and access to mongo
	account = get_account()
	storage = get_storage(namespace='events_log')

	#get data in payload
	start = request.params.get('start', default=None)
	stop = request.params.get('stop', default=None)

	#set returned variable
	output = { } 
	total = None

	#build filter for mongo request (put advanced option in it)
	mfilter = {'$and':[{'crecord_type':'event'}]}

	if start :
		mfilter['and'].append({'timestamp': {'$gt': start}})
	if stop :
		mfilter['and'].append({'timestamp': {'$lt': stop}})
	
	#clean mfilter
	mfilter = clean_mfilter(mfilter)

	records = storage.find(mfilter, account=account)

	tmp_output = {}
	for record in records:
		dump = record.dump()
		if not tmp_output.has_key(dump['event_id']) :
			tmp_output[dump['event_id']] = []
		tmp_output[dump['event_id']].append(dump)

	events = {} 
	for tag in tmp_output.keys() :
		prev_id = None
		for record in tmp_output[tag] :
			item = {}
			if ( prev_id != None ) :
				events[prev_id]['end'] = datetime.fromtimestamp(record['timestamp']).isoformat()#strftime('%Y-%m-%d %H:%M:%S')
			item['title'] = str(dump['component'])
			if 'long_output' in record :
				item['description'] = record['long_output']
			item['start'] = datetime.fromtimestamp(record['timestamp']).isoformat()#strftime('%Y-%m-%d %H:%M:%S')
			item['durationEvent'] = True
			events[record['_id']] = item
			prev_id = record['_id'] 
	output['events'] = events.values()
	output['dateTimeFormat'] = 'iso8601'
	output['wiki-url'] = "http://simile.mit.edu/shelf/"
	output['wiki-section'] = "Simile JFK Timeline"
	#return the total, the success of function and the requested data in output
	return output
