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


@delete('/ui/view/:name')
def tree_delete(name=None):
	namespace='object'
	account = get_account()
	storage = get_storage(namespace=namespace, account=account, logging_level=logging.DEBUG)
	
	record = storage.get(name, account=account)
	
	if isinstance(record, crecord):
		if len(record.children) == 0:
			if record.check_write(account=account):
				#remove record from its parent child list
				for parent in record.parent:
					parent_rec = storage.get(parent, account=account)
					parent_rec.remove_children(record )
					if parent_rec.check_write(account=account):
						storage.put(parent_rec,account=account)
					else:
						logger.debug('Access Denied')
						return HTTPError(403, "Access Denied")
			
				try:
					storage.remove(record, account=account)
				except Exception, err:
					logger.error(err)
					return HTTPError(404, 'Error while removing: %s' % err)
			else:
				logger.debug('Access Denied')
				return HTTPError(403, "Access Denied")
				
		else:
			logger.warning('This record have children, remove those child before')

	
	
@post('/ui/view',checkAuthPlugin={'authorized_grp':group_managing_access})
@post('/ui/view/:name',checkAuthPlugin={'authorized_grp':group_managing_access})
def tree_update(name='None'):
	namespace = 'object'
	account = get_account()

	storage = get_storage(namespace=namespace, account=account)
	
	data = json.loads(request.body.readline())
	if isinstance(data,list):
		data = data[0]
	
	logger.debug('Tree update %s' % data['_id'])

	try:
		logger.debug(' + Get future parent record')
		record_parent = storage.get(data['parentId'], account=account)
	except:
		return HTTPError(403, "You don't have right on the parent record: %s")
	
	try:
		logger.debug(' + Get the children record')
		record_child = storage.get(data['_id'], account=account)
	except:
		logger.debug('   + Children not found')
		record_child = None

	#test if the record exist
	if isinstance(record_child, crecord):
		#check write rights
		if record_parent.check_write(account=account) and record_child.check_write(account=account):
			#if parents are really different
			if not (data['parentId'] in record_child.parent):
				logger.debug(' + Update relations')	
				record_parent_old = storage.get(record_child.parent[0], account=account)
				logger.debug('   + Remove children %s from %s' % (record_child._id, record_parent_old._id))
				record_parent_old.remove_children(record_child)
				
				logger.debug('   + Add children %s to %s' % (record_child._id, record_parent._id))
				record_parent.add_children(record_child)
				
				logger.debug('   + Updating all records')
				storage.put([record_child, record_parent, record_parent_old],account=account)
				
			elif (record_child.name != data['crecord_name']):
				logger.debug(' + Rename record')	
				logger.debug('   + old name : %s' % record_child.name)
				logger.debug('   + new name : %s' % data['crecord_name'])
				record_child.name = data['crecord_name']
				storage.put(record_child,account=account)
			
			else :
				logger.debug(' + Records are same, nothing to do')
		else:
			logger.debug('Access Denied')
			return HTTPError(403, "Access Denied")
			
	else:
		#add new view/folder
		parentNode = storage.get(data['parentId'], account=account)
		#test rights
		if isinstance(parentNode, crecord):
			if parentNode.check_write(account=account):
				if data['leaf'] == True:
					logger.debug('record is a leaf, add the new view')
					record = crecord({'leaf':True,'_id':data['id'],'items':data['items']},type='view',name=data['crecord_name'],account=account)
				else:
					logger.debug('record is a directory, add it')
					record = crecord({'_id':data['id']},type='view_directory',name=data['crecord_name'],account=account)
				
				#must save before add children
				record.chown(account._id)
				record.chgrp(account.group)
				record.chmod('g+w')
				record.chmod('g+r')
				
				storage.put(record,account=account)
				
				parentNode.add_children(record)
				
				
				
				storage.put([record,parentNode],account=account)
			else:
				logger.debug('Access Denied')
				return HTTPError(403, "Access Denied")

		else :
			logger.error('ParentNode doesn\'t exist')

@get('/ui/view/exist/:name')
def check_exist(name=None):
	namespace = 'object'
	account = get_account()
	storage = get_storage(namespace=namespace, account=account)
	
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
	namespace = 'object'
	account = get_account()
	storage = get_storage(namespace=namespace, account=account)
	
	#export_fields = ['']
	
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
