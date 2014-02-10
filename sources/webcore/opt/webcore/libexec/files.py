#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

import sys
import os
import clogging
import json
import gevent

import bottle
from bottle import route, get, delete, put, request
from bottle import HTTPError, post, static_file, response

# GridFS
from pymongo import Connection
import gridfs

# Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord

from cfile import cfile
from cfile import get_cfile
from cfile import namespace

# Import protection function
from libexec.auth import get_account

logger = clogging.getLogger()

# Defines allowed mime types
_allowed_mimetypes = {
	'application/pdf': ['pdf'],
	'image/gif': ['gif'],
	'image/jpeg': ['jpeg', 'jpg'],
	'image/png': ['png'],
	'video/ogg': ['ogg']
}

def get_reversed_http_mimetypes():
	result = {}
	for key, values in _allowed_mimetypes.items():
		for value in values:
			result[value] = key
	return result
allowed_mimetypes = get_reversed_http_mimetypes()
# Max allowed upload file in megabytes
max_size = 5

#########################################################################

@get('/files/:metaId')
@get('/file/:metaId')
@get('/files')
def files(metaId=None):
	
	# Arg option for attachement http option (default: True)
	as_attachment = True
	if request.params.get('as_attachment') in ['false', 'False', '0']:
		as_attachment = False

	if metaId:
		account = get_account()
		storage = get_storage(account=account, namespace=namespace)
	
		logger.debug("Get file '%s'" % metaId)
		
		try:
			rfile = get_cfile(metaId, storage)
		except Exception as err:
			logger.error('File not found: %s' % err)
			return HTTPError(404, "File not found")
		
		file_name = rfile.data['file_name']
		content_type = rfile.data['content_type'] or 'application/octet-stream'
		
		logger.debug(" + File name:    %s" % file_name)
		logger.debug(" + Content type: %s" % content_type)
		logger.debug(" + Bin Id:       %s" % rfile.get_binary_id())
		
		try:
			data = rfile.get()
		except Exception as err:
			logger.error('Error while file fetching: %s' % err)

		if data:
			if as_attachment:
				try:
					response.headers['Content-Disposition'] = 'attachment; filename="%s"' % file_name.encode("utf8")
				except Exception as err:
					logger.error(err)
					return HTTPError(500, "Impossible to encode file_name.")

			response.headers['Content-Type'] = content_type
			try:
				return data
			except Exception as err:
				logger.error(err)
		else:
			logger.error('File not found in gridfs')
			return HTTPError(404, "File not found")
	else:
		return list_files()


@post('/file')
@post('/files')
def add_file():
	# Must be set as text/html cause of extjs upload file method
	# http://docs.sencha.com/extjs/4.0.7/#!/api/Ext.form.Basic-method-hasUpload
	# A json in a string will be return to avoid Bottle to automatically set header to json
	response.headers['Content-Type'] = 'text/html'

	data = request.files['file-path']

	if data.filename and data.file:
		if allowed_mimetypes.get(data.filename.split('.')[-1], False):
			content_type = allowed_mimetypes[data.filename.split('.')[-1]]
			account = get_account()
			storage = get_storage(account=account, namespace=namespace)
			cfile_record = cfile(storage=storage)
			cfile_record.put_data(data.file.read(), file_name=data.filename, content_type=content_type)
			try:
				file_id = storage.put(cfile_record)
				data = {'success': True, 'data': {'code': 200, 'message': 'File uploaded', 'filename': data.filename, 'file_id': str(file_id)}}
			except Exception as err:
				data = {'success': False, 'data': {'code': 500, 'message': err}}
		else:
			data = {'success': False, 'data': {'code': 415, 'message': 'Unsupported Media Type'}}
	else:
		data = {'success': False, 'data': {'code': 400, 'message': 'Bad request'}}

	return json.dumps(data)


@put('/files')
@put('/files/:metaId')
def update_file(metaId=None):	
	data = json.loads(request.body.readline())
	
	if not metaId:
		metaId = data.get('id', None)
	file_name = data.get('file_name', None)
	
	logger.debug("Update file")
	
	if not metaId:
		logger.error('No file Id specified')
		return HTTPError(405, " No file Id specified")
	
	if not file_name:
		logger.error('No file_name specified')
		return HTTPError(405, " No file_name specified")
		
	logger.debug(" + metaId: %s" % metaId)
	logger.debug(" + file name: %s" % file_name)
		
	account = get_account()
	storage = get_storage(account=account, namespace=namespace)
		
	try:
		record = storage.get(metaId)
		if record:
			record.data['file_name'] = file_name
			record.name = file_name
			storage.put(record)
					
	except Exception, err:
		logger.error("Error when updating report %s: %s" % (metaId,err))
		return HTTPError(500, "Failed to update report")


@delete('/files/:metaId')
@delete('/files')
def delete_file(metaId=None):
	account = get_account()
	storage = get_storage(account=account, namespace=namespace)
	
	rfiles = []
	data = request.body.readline()
	items = []
	try:
		items = json.loads(data)
	except:
		logger.warning('Invalid data in request payload')	
		
	## Only accept list for multiremove
	if not isinstance(items, list):
		items = []
	
	if metaId and not items:
		rfile = get_cfile(metaId, storage)
		rfiles.append(rfile)
	else:
		logger.debug('Multi-remove: %s' % data)
			
		for item in items:
			rfile = get_cfile(item['id'], storage)
			rfiles.append(rfile)
						
	logger.debug('Remove %s files' % len(rfiles))
	try:
		for rfile in rfiles:
			rfile.remove()
		return {'total': len(rfiles),"data": [] ,"success":True}
	except:
		logger.error('Failed to remove file')
		return HTTPError(500, "Failed to remove file")

		
		
def list_files():
	limit		= int(request.params.get('limit', default=20))
	start		= int(request.params.get('start', default=0))
	sort		= request.params.get('sort', default=None)
	filter		= request.params.get('filter', default={})
	
	###########account and storage
	account = get_account()
	storage = get_storage(account=account, namespace=namespace)

	logger.debug("List files")
		
	###########load filter
	if filter:
		try:
			filter = json.loads(filter)
		except Exception, err:
			logger.error("Filter decode: %s" % err)
			filter = {}
			
	if isinstance(filter, list):
		if len(filter) > 0:
			filter = filter[0]
		else:
			logger.error(" + Invalid filter format")
			filter = {}
	
	msort = []
	if sort:
		sort = json.loads(sort)
		for item in sort:
			direction = 1
			if str(item['direction']) == "DESC":
				direction = -1

			msort.append((str(item['property']), direction))
	
	
	###########search
	try:
		records = storage.find(filter, sort=msort,limit=limit, offset=start,account=account)
		total = storage.count(filter, account=account)
	except Exception, err:
		logger.error('Error while fetching records: %s' % err)
		return HTTPError(500, "Error while fetching records: %s" % err)
	
	data = []
	
	for record in records:
		data.append(cfile(record=record).dump(json=True))
		
	return {'total': total, 'success': True, 'data': data}
