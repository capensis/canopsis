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

from caccount import caccount
from cstorage import cstorage
from crecord import crecord

##set root account
root = caccount(user="root", group="root")

logger = None

def init():
	storage = cstorage(account=root)
	
	namespaces = ['files', 'binaries.files', 'binaries.chunks']
	
	for namespace in namespaces:
		logger.info(" + Drop '%s' collection" % namespace)
		storage.drop_namespace(namespace)

def update():
	#storage = cstorage(account=root)
	
	#logger.info(" + Update: 'daily -> 201205'")
	#collections = storage.db.collection_names()
	#if 'reports' in collections:
	#	logger.info("   + Move records from 'reports' to 'files' ...")
	#	records = storage.find({}, account=root, namespace='reports')
	#	for record in records:
	#		logger.info("     - %s (%s)" % (record._id, record.data['file_name']))
	#		storage.put(record, account=root, namespace='files')
	#		
	#	storage.drop_namespace('reports')
	update_for_new_rights()
	
	pass
		
def update_for_new_rights():
	#update briefcase elements
	storage = cstorage(namespace='files',account=root)
	
	dump = storage.find({})

	for record in dump:
		if record.owner.find('account.') == -1:
			record.owner = 'account.%s' % record.owner
		if record.group.find('group.') == -1:
			record.group = 'group.%s' % record.group
			
	storage.put(dump)
