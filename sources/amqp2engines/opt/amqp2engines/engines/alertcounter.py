#!/usr/bin/env python
#--------------------------------
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

import logging
import time

from copy import deepcopy

from cengine import cengine
import cevent

from cstorage import get_storage
from caccount import caccount
import pyperfstore2

NAME="alertcounter"
INTERNAL_COMPONENT = '__canopsis__'
PROC_CRITICAL = 'PROC_CRITICAL'
PROC_WARNING = 'PROC_WARNING'

class engine(cengine):
	def __init__(self, *args, **kargs):
		super(engine, self).__init__(name=NAME, *args, **kargs)

	def pre_run(self):
		self.listened_event_type = ['check','selector','eue','sla', 'log']
		self.manager = pyperfstore2.manager()

		# Get SLA configuration
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.entities = self.storage.get_backend('entities')

		self.beat()

	def load_macro(self):
		self.logger.debug('Load record for macros')

		self.mCrit = PROC_CRITICAL
		self.mWarn = PROC_WARNING

		record = self.storage.find_one({'crecord_type': 'sla', 'objclass': 'macro'})

		if record:
			self.mCrit = record.data['mCrit']
			self.mWarn = record.data['mWarn']

	def load_crits(self):
		self.logger.debug('Load records for criticalness')

		self.crits = {}

		records = self.storage.find({'crecord_type': 'sla', 'objclass': 'crit'})

		for record in records:
			self.crits[record.data['crit']] = record.data['delay']

		self.selectors_name = []
		self.last_resolv = 0

	def beat(self):
		self.load_macro()
		self.load_crits()

	def perfdata_key(self, meta):
		if 're' in meta and meta['re']:
			return '{0}{1}{2}'.format(meta['co'], meta['re'], meta['me'])

		else:
			return '{0}{1}'.format(meta['co'], meta['me'])

	def increment_counter(self, key, meta, value):
		self.logger.debug("Increment {0}: {1}".format(key, value))
		self.manager.push(name=key, value=value, meta_data=meta)

	def update_global_counter(self, event):
		# Update global counter
		event = cevent.forger(
			connector = "cengine",
			connector_name = NAME,
			event_type = "check",
			source_type = "component",
			component = INTERNAL_COMPONENT,

			state = event['state'],
			state_type = event['state_type'],
			component_problem = event['component_problem']
		)

		self.amqp.publish(event, cevent.get_routingkey(event), self.amqp.exchange_name_events)

		event['source_type'] = 'resource'

		for hostgroup in event.get('hostgroups', []):
			event['resource'] = hostgroup

			self.amqp.publish(event, cevent.get_routingkey(event), self.amqp.exchange_name_events)


	def count_sla(slatype, delay):
		meta_data = {'type': 'COUNTER', 'co': INTERNAL_COMPONENT }
		now = int(time.time())

		if delay < (now - event['last_state_change']):
			ack = self.entities.find_one({
				'type': 'ack',
				'timestamp': {
					'$gt': event['last_state_change'],
					'$lt': event['previous_state']
				}
			})

			meta_data['me'] = 'cps_sla_{0}_{1}_{2}'.format(
				slatype,
				warn.lower(),
				'nok' if ack else 'out'
			)

		else:
			meta_data['me'] = 'cps_sla_{0}_{1}_ok'.format(slatype, warn.lower())

		key = self.perfdata_key(meta_data)
		self.increment_counter(key, meta_data, 1)

	def count_by_crits(self, event):
		if event['state'] == 0 and event.get('state_type', 1) == 1:
			warn = event.get(self.mWarn, None)
			crit = event.get(self.mCrit, None)

			if warn and warn in self.crits and event['previous_state'] == 1:
				self.count_sla('warn', self.crits[warn])

			elif crit and crit in self.crits and event['previous_state'] == 2:
				self.count_sla('crit', self.crits[crit])

	def count_alert(self, event):
		component = event['component']
		resource = event.get('resource', None)
		tags = event.get('tags', [])

		# Update cps_statechange{,_0,_1,_2,_3} for component/resource

		meta_data = {
			'type': 'COUNTER',
			'co': component,
			're': resource,
			'tg': tags
		}

		meta_data['me'] = "cps_statechange"
		key = self.perfdata_key(meta_data)

		self.increment_counter(key, meta_data, value)

		meta_data['me'] = "cps_statechange_nok"
		key = self.perfdata_key(meta_data)

		cvalue = value if state != 0 else 0

		self.increment_counter(key, meta_data, cvalue)

		for cstate in range(3):
			cvalue = value if cstate == state else 0

			meta_data['me'] = "cps_statechange_{0}".format(cstate)
			key = self.perfdata_key(meta_data)

			self.increment_counter(key, meta_data, cvalue)

		# Update cps_statechange_{hard,soft}

		for cstate_type in range(1):
			cvalue = value if cstate_type == state_type else 0

			meta_data['me'] = "cps_statechange_{0}".format(
				'hard' if cstate_type == 1 else 'soft'
			)

			key = self.perfdata_key(meta_data)

			self.increment_counter(key, meta_data, cvalue)

		# Update cps_statechange_{component,resource,resource_by_component}

		if component == INTERNAL_COMPONENT:
			for cevtype in ['component', 'resource', 'resource_by_component']:
				cvalue = 0

				if cevtype == 'component' and not resource:
					cvalue = value

				elif cevtype == 'resource' and resource and not cmp_problem:
					cvalue = value

				elif cevtype == 'resource_by_component' and resource and cmp_problem:
					cvalue = value

				meta_data['me'] = "cps_statechange_{0}".format(cevtype)
				key = self.perfdata_key(meta_data)

				self.increment_counter(key, meta_data, cvalue)

		# Update cps_alerts_not_ack

		if state != 0:
			metric = "cps_alerts_not_ack"
			meta_data['me'] = metric

			self.logger.debug("Increment %s: %s: %s" % (name, metric, value))
			self.manager.push(name="%s%s" % (name, metric), value=value, meta_data=meta_data)

	def resolve_selectors_name(self):
		if int(time.time()) > (self.last_resolv + 60):
			records = self.storage.find(mfilter={'crecord_type': 'selector'}, mfields=['crecord_name'])

			self.selectors_name = [record['crecord_name'] for record in records]

			self.last_resolv = int(time.time())

		return self.selectors_name

	def count_by_tags(self, event, value):
		if event['event_type'] != 'selector':
			tags = event.get('tags', [])
			tags = [tag for tag in tags if tag in self.resolve_selectors_name()]

			for tag in tags:
				self.logger.debug("Increment Tag: '%s'" % tag)
				tagevent = deepcopy(event)
				tagevent['component'] = tag
				tagevent['resource'] = 'selector'

				self.count_alert(tagevent, value)

	def work(self, event, *args, **kargs):
		validation = event['event_type'] in self.listened_event_type
		validation = validation and event.get('state_type', 1) == 1
		validation = validation and event['component'] != 'derogation'

		if validation:
			self.update_global_counter(event)
			self.count_by_crits(event)

			# By name
			self.count_alert(event, 1)

			# By tags (selector)
			self.count_by_tags(event, 1)

			# By crits
			self.count_by_crits(event, 1)

		return event
