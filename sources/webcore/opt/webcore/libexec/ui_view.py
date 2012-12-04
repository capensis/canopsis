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
from libexec.auth import check_auth, get_account, check_group_rights

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
@delete('/ui/view/:name')
def tree_delete(name=None):
	logger.debug('DELETE:')
	account = get_account()
	storage = get_storage(namespace='object', account=account, logging_level=logging.DEBUG)
	
	if not name:
		name = []
		data = json.loads(request.body.readline())
		if not isinstance(data,list):
			data = [data]
		for view in data:
			name.append(view['_id'])
			

	output = {}
	
	for _id in name:
		try:
			record = storage.get(_id, account=account)
		except Exception, err:
			logger.info(' + Record not found: %s' %_id)
			output[record._id] = {'success':False,'output':'Record not found'}
			record = None
			
		if record:
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
	
	return {"total": len(name), "success": True, "data": output}

@post('/ui/view',checkAuthPlugin={'authorized_grp':group_managing_access})
def tree_add():
	logger.debug('POST:')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	data = json.loads(request.body.readline())
	if not isinstance(data,list):
		data = [data]
		
	output = {}

	for view in data:
		view_name = view.get('crecord_name',view['_id'])
		
		try:
			logger.debug(' + Get future parent record')
			record_parent = storage.get(view['parentId'], account=account)
		except:
			logger.info("You don't have right on the parent record: %s" % view['parentId'])
			output[view_name] = {'success':False,'output':"You don't have right on the parent record"}
			record_parent = None
		
		try:
			record_child = storage.get(view['_id'], account=account)
			logger.debug(' + Children found %s' % view['_id'])
			output[view['_id']] = {'success':False,'output':"Record id already exist"}
		except:
			logger.debug(' + Children not found')
			record_child = None

		if record_parent and not record_child:
			if record_parent.check_write(account=account):
				try:
					if view['leaf'] == True:
						logger.debug('record is a leaf, add the new view')
						record = crecord({'leaf':True,'_id':view['id'],'items':view['items']},type='view',name=view['crecord_name'],account=account)
					else:
						logger.debug('record is a directory, add it')
						record = crecord({'_id':view['id']},type='view_directory',name=view['crecord_name'],account=account)
				except Exception, err:
					logger.info('Error while building view/directory crecord : %s' % err)
					output[view_name] = {'success':False,'output':"Error while building crecord: %s" % err}
					record = None
					
				if isinstance(record,crecord):
					record.chown(account._id)
					record.chgrp(account.group)
					record.chmod('g+w')
					record.chmod('g+r')
					
					storage.put(record,account=account)
					record_parent.add_children(record)

					storage.put([record,record_parent],account=account)
					output[view_name] = {'success':True,'output':''}
			else:
				logger.info('Access Denied')
				output[view_name] = {'success':False,'output':"No rights on this record"}
		else:
			logger.error("Parent doesn't exists or view/directory already exists for %s" % view_name)
			
	return {"total": len(data), "success": True, "data": output}

@put('/ui/view',checkAuthPlugin={'authorized_grp':group_managing_access})
def update_view_relatives():
	logger.debug('PUT:')
	account = get_account()
	storage = get_storage(namespace='object', account=account)
	
	data = json.loads(request.body.readline())
	if not isinstance(data,list):
		data = [data]
	
	output={}
	
	for view in data:
		view_name = view.get('crecord_name',view['_id'])
		_id = view.get('_id',None)
		parent_id = view.get('parentId',None)
		
		logger.debug(' + View to update is %s' % _id)
		logger.debug(' + Parent is %s' % parent_id)
		
		#get records
		record_child=None
		try:
			record_child = storage.get(_id, account=account)
		except Exception, err:
			logger.error(" + Record to update wasn't found")
			output[view_name] = {'success':False,'output':str(err)}

		
		#check action to do
		if record_child:
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
