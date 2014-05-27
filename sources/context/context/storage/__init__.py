#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from cmongo.storage.timed import TimedStorage


_STORES_BY_COLLECTION = dict()

DEFAULT_TABLE = 'context'


class TimedStorage(TimedStorage):
	"""
	Manage access to context (connector, component, resource) entities
	and context data (metric, downtime, etc.) related to context entities.
	"""

	DEFAULT_CONFIGURATION_FILE = '~/etc/context.conf'

	class Index:
		"""
		Provide keys and index to database.
		"""

		TYPE = 'type'
		NAME = 'NAME'
		CONNECTOR = 'connector'
		COMPONENT = 'component'
		RESOURCE = 'resource'
		ID = 'id'
		NODEID = 'nodeid'
		_ID = '_id'

		NAME_BY_TYPE = [(TYPE, 1), (NAME, 1)]
		DOWNTIME = [
			(TYPE, 1), (CONNECTOR, 1), (COMPONENT, 1), (RESOURCE, 1)]
		CONTEXT = [
			(TYPE, 1), (CONNECTOR, 1), (COMPONENT, 1), (RESOURCE, 1), (NAME, 1)]

	def _get_indexes():

		result = [
			[('_id')],
			TimedStorage.Index.NAME_BY_TYPE,
			TimedStorage.Index.DOWNTIME,
			TimedStorage.Index.CONTEXT
		]

		return result

	def __init__(self, backend=DEFAULT_TABLE,
		configuration_file=DEFAULT_CONFIGURATION_FILE, *args, **kwargs):

		super(TimedStorage, self).__init__(
			backend=backend, configuration_file=configuration_file,
			*args, **kwargs)

	def entity_by_id(self, _id):

		document = {
			TimedStorage.Index._ID: _id}

		cursor = self._find(document=document)

		result = None

		for result in cursor:
			break

		return result

	def entity_by_name(
		self, connector=None, connector_type=None, component=None, resource=None,
		data_context=None, data_context_type=None):

		query = {
			TimedStorage.Index.CONNECTOR: connector,
		}

		query[TimedStorage.Index.NAME] = connector
		query[TimedStorage.Index.TYPE] = TimedStorage.Index.CONNECTOR

		if component is not None:
			query[TimedStorage.Index.CONNECTOR] = connector
			query[TimedStorage.Index.NAME] = component
			query[TimedStorage.Index.TYPE] = TimedStorage.Index.COMPONENT

		if resource is not None:
			query[TimedStorage.Index.COMPONENT] = component
			query[TimedStorage.Index.NAME] = resource
			query[TimedStorage.Index.TYPE] = TimedStorage.Index.RESOURCE

		if data_context is not None:
			query[TimedStorage.Index.COMPONENT] = component
			query[TimedStorage.Index.NAME] = data_context
			query[TimedStorage.Index.TYPE] = data_context_type

		cursor = self._find(document=query)

		cursor.hint(TimedStorage.Index.CONTEXT)

		result = None

		for result in cursor:
			break

		return result

	def entities_by_name(self, entity_type, entity_name):

		query = {
			TimedStorage.Index.TYPE: entity_type,
			TimedStorage.Index.NAME: entity_name
		}

		cursor = self._find(document=query)

		cursor.hint(TimedStorage.Index.NAME_BY_TYPE)

		return cursor

	def get_entities_from_db(self, mfilter, backend=None):

		cursor = self._find(document=mfilter)

		result = list(cursor)

		for document in result:
			# hack for ObjectId serialization
			document[TimedStorage.Index._ID] = str(document[TimedStorage.Index._ID])

		return result
