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

import clogging

from bottle import get, post, request, HTTPError

from caccount import caccount
from cstorage import get_storage

from libexec.auth import get_account

from bson import json_util
import json


logger = clogging.getLogger()


def get_entities_from_db(mfilter):
	account = get_account()
	storage = get_storage(namespace='entities', account=account)
	backend = storage.get_backend()

	entities = backend.find(mfilter)

	response = {
		'total': entities.count(),
		'data': [],
		'success': True
	}

	for entity in entities:
		tmp = entity
		tmp['_id'] = str(tmp['_id'])

		response['data'].append(tmp)

	return response

@get('/entities/')
def get_all_entities():
	return get_entities_from_db({})

@get('/entities/:etype')
def get_all_typed_entities(etype):
	return get_entities_from_db({
		'type': etype
	})

@get('/entities/:etype/:ename')
def get_specific_entity(etype, ename):
	identifier = 'name'

	# Specific cases :
	if etype == 'downtime':
		identifier = 'id'

	elif etype == 'ack':
		identifier = 'timestamp'

	return get_entities_from_db({
		'type': etype,
		identifier: ename
	})

@post('/entities/')
def get_entities_with_custom_filter():
	try:
		mfilter = request.body.readline()
		mfilter = json.loads(mfilter)['filter']

	except ValueError, err:
		if 'filter' not in request.params:
			return HTTPError(500, json.dumps({
				'total': 0,
				'data': [],
				'success': False,
				'error': str(err)
			}))

		mfilter = request.params['filter']

	return get_entities_from_db(mfilter)