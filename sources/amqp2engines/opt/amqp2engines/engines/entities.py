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

from cengine import cengine
from caccount import caccount
from cstorage import get_storage
import cevent
import time
import md5

NAME="entities"

class engine(cengine):
	def __init__(self, name=NAME, *args, **kwargs):
		super(engine, self).__init__(name=name, *args, **kwargs)

		self.account = caccount(user='root', group='root')
		self.storage = get_storage(namespace='entities', logging_level=self.logging_level, account=self.account)
		self.backend = self.storage.get_backend()

	def work(self, event, *args, **kwargs):
		mCrit = 'PROC_CRITICAL'
		mWarn = 'PROC_WARNING'

		record = self.storage.find_one(
			mfilter = {
				'crecord_type': 'sla',
				'objclass': 'macro'
			},
			namespace='object'
		)

		if record:
			mCrit = record.data['mCrit']
			mWarn = record.data['mWarn']

		# Get event informations
		component = event['component']
		resource = event.get('resource', None)
		hostgroups = event.get('hostgroups', [])
		servicegroups = event.get('servicegroups', [])

		# Create Component entity
		data = {
			'type': 'component',
			'name': component,
			'hostgroups': hostgroups,
			'mCrit': None,
			'mWarn': None
		}

		if event['source_type'] == 'component':
			data['mCrit'] = event.get(mCrit, None)
			data['mWarn'] = event.get(mWarn, None)

		self.backend.update({
				'type': 'component',
				'name': component
			},{
				'$set': data
			},
			upsert = True
		)

		# Create Resource entity
		if resource:
			self.backend.update({
					'type': 'resource',
					'name': resource
				},{
					'$set': {
						'type': 'resource',
						'name': resource,
						'component': component,
						'hostgroups': hostgroups,
						'servicegroups': servicegroups,
						'mCrit': event.get(mCrit, None),
						'mWarn': event.get(mWarn, None)
					}
				},
				upsert = True
			)

		# Create Hostgroups entities
		for hostgroup in hostgroups:
			self.backend.update({
					'type': 'hostgroup',
					'name': hostgroup
				},{
					'$set': {
						'type': 'hostgroup',
						'name': hostgroup
					}
				},
				upsert = True
			)

		# Create Servicegroups entities
		for servicegroup in servicegroups:
			self.backend.update({
					'type': 'servicegroup',
					'name': servicegroup
				},{
					'$set': {
						'type': 'servicegroup',
						'name': servicegroup
					}
				},
				upsert = True
			)

		# Create Downtime entity
		if event['event_type'] == 'downtime':
			self.backend.update({
					'type': 'downtime',
					'component': component,
					'resource': resource,
					'id': event['downtime_id']
				},{
					'$set': {
						'type': 'downtime',
						'component': component,
						'resource': resource,
						'id': event['downtime_id'],

						'author': event['author'],
						'comment': event['output'],

						'start': event['start'],
						'end': event['end'],
						'duration': event['duration'],

						'fixed': event['fixed'],
						'entry': event['entry']
					}
				},
				upsert = True
			)

		# Create acknowledgement entity
		elif event['event_type'] == 'ack':
			self.backend.update({
					'type': 'ack',
					'timestamp': event['timestamp'],
					'component': component,
					'resource': resource,
				},{
					'$set': {
						'type': 'ack',
						'timestamp': event['timestamp'],
						'component': component,
						'resource': resource,

						'author': event['author'],
						'comment': event['output'],
					}
				},
				upsert = True
			)

		# Create metrics entities
		for perfdata in event['perf_data_array']:
			nodeid = md5.new()

			nodeid.update(component)

			if resource:
				nodeid.update(resource)

			nodeid.update(perfdata['metric'])
			nodeid = nodeid.hexdigest()

			self.backend.update({
					'type': 'metric',
					'nodeid': nodeid
				},{
					'$set': {
						'type': 'metric',
						'component': component,
						'resource': resource,
						'name': perfdata['metric'],
						'nodeid': nodeid,
						'last': [event['timestamp'], perfdata['value']],
						'min': perfdata.get('min', None),
						'max': perfdata.get('max', None),
						'warn': perfdata.get('warn', None),
						'crit': perfdata.get('crit', None),
						'unit': perfdata.get('unit', None),
						'perftype': perfdata.get('type', 'GAUGE')
					}
				},
				upsert = True
			)

		return event