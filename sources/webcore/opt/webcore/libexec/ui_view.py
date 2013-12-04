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
from bottle import route, get, put, delete, request, HTTPError, post,response

## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord

#import protection function
from libexec.auth import get_account, check_group_rights

logger = logging.getLogger("ui_view")

#group who have right to access 
group_managing_access = ['group.CPS_view_admin','group.CPS_view']
export_fields = ['items','crecord_name']

#########################################################################

@get('/ui/view')
def tree_get():
	namespace = 'object'
	account = get_account()
		
	storage = get_storage(namespace=namespace, account=account, logging_level=logging.DEBUG)
	
	node = request.params.get('node', default= None)
	
	output = []
	total = 0
		
	if node:
		parentNode = storage.get('directory.root', account=account)
		storage.recursive_get(parentNode,account=account)
		output = parentNode.recursive_dump(json=True)
			
	return output


@delete('/ui/view')
@delete('/ui/view/:_id')
def tree_delete(_id=None):
	logger.debug('DELETE:')
	account = get_account()
	storage = get_storage(namespace='object', account=account, logging_level=logging.DEBUG)
	
	ids = []

	try:
		data = json.loads(request.body.readline())
		if not isinstance(data,list):
			data = [data]
		for view in data:
			ids.append(view['_id'])
	except:
		logger.debug('No payload for delete action')
		ids.append(_id)

	output = {}
	
	for _id in ids:
		try:
			record = storage.get(_id, account=account)
		except Exception, err:
			logger.info(' + Record not found: %s' %_id)
			output[_id] = {'success':False, 'output':'Record not found'}
			continue
			
		if len(record.children) == 0:
			if record.check_write(account=account):
				#remove record from its parent child list
				for parent in record.parent:
					parent_rec = storage.get(parent, account=account)
					parent_rec.remove_children(record )
					if parent_rec.check_write(account=account):
						storage.put(parent_rec,account=account)
					else:
						logger.info(' + No sufficient rights on parent %s' % parent_rec.name)
						output[parent_rec.name] = {'success':False,'output':'No right to remove children'}
				#remove the record
				try:
					storage.remove(record, account=account)
					output[record.name] = {'success':True,'output':''}
				except Exception, err:
					logger.error(err)
					output[record.name] = {'success':False,'output':err}
			else:
				logger.info(' + No sufficient rights on %s' % parent_rec.crecord_name)
				output[record.name] = {'success':False,'output':'No sufficient rights'}
		else:
			output[record.name] = {'success':False,'output':'This record have children, remove those child before'}
			logger.warning('This record have children, remove those child before')
	
	return {"total": len(ids), "success": True, "data": output}

@post('/ui/view',checkAuthPlugin={'authorized_grp':group_managing_access})
@post('/ui/view/:_id',checkAuthPlugin={'authorized_grp':group_managing_access})
def tree_add(_id=None):
	logger.debug('POST:')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	data = json.loads(request.body.readline())
		
	output = add_view(data, storage, account)
			
	return {"total": len(data), "success": True, "data": output}

@put('/ui/view',checkAuthPlugin={'authorized_grp':group_managing_access})
@put('/ui/view/:_id',checkAuthPlugin={'authorized_grp':group_managing_access})
def update_view_relatives(_id=None):
	logger.debug('PUT:')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	data = json.loads(request.body.readline())
	
	output = update_view(data, storage, account)

	return {"total": len(data), "success": True, "data": output}


@get('/ui/view/exist/:name')
def check_exist(name=None):
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	mfilter = {'crecord_name':name}
	
	try:
		logger.debug('try to get view')
		record_child = storage.find_one(mfilter=mfilter, account=account)
		if record_child:
			return {"total": 1, "success": True, "data": {'exist' : True}}
		else:
			return {"total": 0, "success": True, "data": {'exist' : False}}
	except Exception,err:
		logger.error('Error while fetching view : %s' % err)
		return {"total": 0, "success": False, "data": {}}

@get('/ui/view/export/:_id')
def exportView(_id=None):
	logger.debug('Prepare to return view json file')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	if not _id:
		_id = request.params.get('_id', default=None)
	
	try:
		logger.debug(' + Try to get view from database')
		record = storage.get(_id, account=account)
		
		logger.debug(' + %s found' % record.name)
		
		response.headers['Content-Disposition'] = 'attachment; filename="%s.json"' % record.name
		response.headers['Content-Type'] = 'application/json'
		
		record.parent = []
		record._id = None
		
		dump = record.dump()
		output = {}
		
		for item in dump:
			if item in export_fields:
				output[item] = dump[item]
		
		return json.dumps(output, sort_keys=True, indent=4)
		
	except Exception,err:
		logger.error(' + Error while fetching view : %s' % err)
		return {"total": 0, "success": False, "data": {}}

@get('/ui/export/object/:_id')
def export_object(_id=None):
	logger.debug('Prepare to return object json file')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	if not _id:
		_id = request.params.get('_id', default=None)
	
	try:
		logger.debug(' + Try to get object from database')
		record = storage.get(_id, account=account)
		
		logger.debug(' + %s found' % record.name)
		
		response.headers['Content-Disposition'] = 'attachment; filename="%s.json"' % record.name
		response.headers['Content-Type'] = 'application/json'

		dump = record.dump()
		del dump['_id']
	
		return json.dumps(dump, sort_keys=True, indent=4)
		
	except Exception,err:
		logger.error(' + Error while fetching object : %s' % err)
		return json.dumps({'error': str(err)}, sort_keys=True, indent=4)

@post('/ui/export/objects')
def export_objects():
	logger.debug('Prepare to return objects json file')
	account = get_account()
	storage = get_storage(namespace='object', account=account)

	ids = request.params.getall('ids')

	try:
		records = []

		for _id in ids:
			logger.debug(' + Try to get object: %s' % _id)
			record = storage.get(_id, account=account)

			logger.debug(' + %s found' % record.name)

			dump = record.dump()
			del dump['_id']

			records.append(dump)

		response.headers['Content-Disposition'] = 'attachment; filename="objects.json"'
		response.headers['Content-Type'] = 'application/json'

		return json.dumps(records, sort_keys=True, indent=4)

	except Exception, err:
		logger.error(' + Error while fetching objects : %s' % err)
		return json.dumps({'error': str(err)}, sort_keys=True, indent=4)

def add_view(views, storage, account):
	if not isinstance(views, list):
		views = [ views ]

	logger.debug('Create views:')
	output={}

	for view in views:
		view_name = view.get('crecord_name', view['_id'])
		
		record_parent = None
		try:
			logger.debug(' + Get future parent record')
			record_parent = storage.get(view['parentId'], account=account)
		except:
			logger.info("You don't have right on the parent record: %s" % view['parentId'])
			output[view_name] = {'success':False,'output':"You don't have right on the parent record"}
		
		record_child = None
		try:
			record_child = storage.get(view['_id'], account=account)
			logger.debug(' + View already exist %s' % view['_id'])
			output[view['_id']] = {'success':False,'output':"View already exist"}
		except:
			logger.debug(' + View not found')
			
		if record_parent and not record_child:
			if record_parent.check_write(account=account):
				try:
					if view['leaf'] == True:
						logger.debug('record is a leaf, add the new view')
						record = crecord({'leaf':True,'_id':view['_id'],'items':view['items']},type='view',name=view['crecord_name'],account=account)
					else:
						logger.debug('record is a directory, add it')
						record = crecord({'_id':view['_id']},type='view_directory',name=view['crecord_name'],account=account)
				except Exception, err:
					logger.info('Error while building view/directory crecord : %s' % err)
					output[view_name] = {'success':False,'output':"Error while building crecord: %s" % err}
					record = None
					
				if isinstance(record,crecord):
					record.chown(account._id)
					record.chgrp(account.group)

					record.access_group = []
					
					storage.put(record,account=account)
					record_parent.add_children(record)

					storage.put([record,record_parent],account=account)
					output[view_name] = {'success':True,'output':''}
			else:
				logger.info('Access Denied')
				output[view_name] = {'success':False,'output':"No rights on this record"}
		else:
			logger.error("Parent doesn't exists or view/directory already exists for %s" % view_name)

	return output

def update_view(views, storage, account):
	if not isinstance(views, list):
		views = [ views ]

	logger.debug('Update views:')
	output={}

	for view in views:
		view_name = view.get('crecord_name',view['_id'])
		_id = view.get('_id', None)
		parent_id = view.get('parentId', None)
		
		logger.debug(' + View to update is %s' % _id)
		logger.debug(' + Parent is %s' % parent_id)
		
		#get records
		record_child=None
		try:
			record_child = storage.get(_id, account=account)
		except Exception, err:
			logger.error(" + View to update wasn't found, try to add view")
			#output[view_name] = {'success':False,'output':str(err)}
			#### TODO: Check why update was called
			output[view_name] = add_view(view, storage, account)[view_name]
		
		#check action to do
		if record_child:
			# Change tree
			if parent_id and not parent_id in record_child.parent:
				try:
					record_parent = storage.get(parent_id, account=account)
					logger.debug(' + Update relations')	
					record_parent_old = storage.get(record_child.parent[0], account=account)
					logger.debug('   + Remove children %s from %s' % (record_child._id, record_parent_old._id))
					
					if not record_parent.check_write(account=account) or not record_parent_old.check_write(account=account):
						raise Exception('No rights on parent record')
					
					record_parent_old.remove_children(record_child)
					
					logger.debug('   + Add children %s to %s' % (record_child._id, record_parent._id))
					record_parent.add_children(record_child)
					
					logger.debug('   + Updating all records')
		
					storage.put([record_parent,record_child,record_parent_old],account=account)
				except Exception, err:
					output[view_name] = {'success':False,'output':str(err)}
					logger.error(err)

			# Rename view
			elif view['crecord_name'] and record_child.name != view['crecord_name']:
				logger.debug(' + Rename record')	
				logger.debug('   + old name : %s' % record_child.name)
				logger.debug('   + new name : %s' % view['crecord_name'])
				record_child.name = view['crecord_name']
				storage.put(record_child,account=account)
			
			else :
				logger.debug(' + Records are same, nothing to do')
			
			if not view_name in output:
				output[view_name] = {'success':True,'output':""}

	return output
