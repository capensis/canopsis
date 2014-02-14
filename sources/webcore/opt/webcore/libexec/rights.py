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

import sys, os, clogging, json

import bottle
from bottle import route, get, put, delete, request, HTTPError, post

## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord

#import protection function
from libexec.auth import get_account

logger = clogging.getLogger()

#########################################################################

@put('/rights/:namespace/:crecord_id')
def change_rights(namespace,crecord_id=None):
	account = get_account()
	storage = get_storage(namespace=namespace, account=account)
	
	#get put data
	aaa_owner = request.params.get('aaa_owner', default=None)
	aaa_group = request.params.get('aaa_group', default=None)
	aaa_access_owner = request.params.get('aaa_access_owner', default=None)
	aaa_access_group = request.params.get('aaa_access_group', default=None)
	aaa_access_other = request.params.get('aaa_access_other', default=None)
	
	
	if(crecord_id != None):
		record = storage.get(crecord_id, account=account)
		
	if isinstance(record, crecord):
		logger.debug('record found, changing rights/owner')
		#change owner and group
		if aaa_owner is not None:
			record.chown(aaa_owner)
		if aaa_group is not None:
			record.chgrp(aaa_group)
		
		#change rights
		if aaa_access_owner is not None:
			record.access_owner = json.loads(aaa_access_owner)
		if aaa_access_group is not None:
			record.access_group = json.loads(aaa_access_group)
		if aaa_access_other is not None:
			record.access_other = json.loads(aaa_access_other)
			
		#logger.debug(json.dumps(record.dump(json=True), sort_keys=True, indent=4))
		try:
			storage.put(record,account=account)
		except:
			logger.error('Access denied')
			return HTTPError(403, "Access denied")
	
	else:
		logger.warning('The record doesn\'t exist')
	
