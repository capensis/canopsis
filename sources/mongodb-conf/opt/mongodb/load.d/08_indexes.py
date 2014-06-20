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
from cstorage import get_storage

logger = None

##set root account
root = caccount(user="root", group="root")
storage = get_storage(account=root, namespace='object')

INDEXES = {
	'object': [
		[('crecord_type', 1)]
	],
	'perfdata2': [
		[('co', 1), ('re', 1), ('me', 1)],
		[('re', 1), ('me', 1)],
		[('co', 1), ('me', 1)],
		[('me', 1)],
		[('tg', 1)]
	],
	'perfdata2_daily': [
		[('insert_date', 1)]
	],
	'events': [
		[
			('connector_name', 1),
			('event_type', 1),
			('component', 1),
			('resource', 1),
			('state_type', 1),
			('state', 1)
		],[
			('source_type', 1),
			('tags', 1)
		],[
			('event_type', 1),
			('component', 1),
			('resource', 1)
		],[
			('event_type', 1),
			('resource', 1)
		],[
			('event_type', 1)
		]
	],
	'events_log': [
		[
			('connector_name', 1),
			('event_type', 1),
			('component', 1),
			('resource', 1),
			('state_type', 1),
			('state', 1)
		],[
			('source_type', 1),
			('tags', 1)
		],[
			('event_type', 1),
			('component', 1),
			('resource', 1)
		],[
			('event_type', 1),
			('resource', 1)
		],[
			('event_type', 1)
		],[
			('state_type', 1)
		],[
			('tags', 1)
		],[
			('referer', 1)
		]
	],
	'entities': [
		[('type', 1)],
		[('type', 1), ('name', 1)],
		[('type', 1), ('component', 1), ('name', 1)],
		[('type', 1), ('component', 1), ('resource', 1), ('id', 1)],
		[('type', 1), ('component', 1), ('resource', 1)],
		[('type', 1), ('nodeid', 1)]
	]
}

def init():
	for collection in INDEXES:
		logger.info(' + Create indexes for collection {0}'.format(collection))
		col = storage.get_backend(collection)
		col.drop_indexes()

		for index in INDEXES[collection]:
			col.ensure_index(index)

def update():
	init()
