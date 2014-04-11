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

import time

from copy import deepcopy

from cengine import cengine
import cevent

from cstorage import get_storage
from caccount import caccount
import pyperfstore2

import logging

NAME="alertcounter"
INTERNAL_COMPONENT = '__canopsis__'
MACRO = 'CAN_PRIORITY'


class engine(cengine):
	def __init__(self, *args, **kargs):
		super(engine, self).__init__(name=NAME, *args, **kargs)

	def pre_run(self):
		self.listened_event_type = ['check','selector','eue','sla', 'log']
		self.manager = pyperfstore2.manager()

		# Get SLA configuration
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.entities = self.storage.get_backend('entities')
		self.objects_backend = self.storage.get_backend('object')

		self.selectors_name = []
		self.last_resolv = 0
		self.beat()

	def load_macro(self):
		self.logger.debug('Load record for macros')

		self.MACRO = MACRO


		record = self.storage.get_backend('object').find_one({'crecord_type': 'slamacros'})

		if record and 'macro' in record:
			self.MACRO = record['macro']

	def load_crits(self):
		self.logger.debug('Load records for criticalness')

		self.crits = {}

		records = self.storage.find({'crecord_type': 'slacrit'})

		for record in records:
			self.crits[record.data['crit']] = record.data['delay']

	def reload_ack_comments(self):

		# reload comment for ack comparison
		self.comments = {}
		query = self.objects_backend.find({'crecord_type': 'comment', 'referer_event_rks': {'$exists': True} }, {'referer_event_rks':1})
		for comment in query:
			for rk in comment['referer_event_rks']:
				self.comments[rk['rk']] = 1
		self.logger.debug('loaded %s referer key comments' % len(self.comments))

	def beat(self):
		self.load_macro()
		self.load_crits()
		self.reload_ack_comments()

	def perfdata_key(self, meta):
		if 're' in meta and meta['re']:
			return u'{0}{1}{2}'.format(meta['co'], meta['re'], meta['me'])

		else:
			return u'{0}{1}'.format(meta['co'], meta['me'])

	def increment_counter(self, meta, value):
		key = self.perfdata_key(meta)
		self.logger.debug(u"Increment {0}: {1}".format(key, value))
		self.logger.debug(str(meta))
		self.manager.push(name=key, value=value, meta_data=meta)

	def update_global_counter(self, event):
		# Comment action (ensure the component exists in database)
		logevent = cevent.forger(
			connector = 'cengine',
			connector_name = NAME,
			event_type = 'log',
			source_type = 'component',
			component = INTERNAL_COMPONENT,

			output = 'Updating global counter'
		)

		self.amqp.publish(logevent, cevent.get_routingkey(logevent), self.amqp.exchange_name_events)

		# Update counter
		new_event = deepcopy(event)
		new_event['connector']      = 'cengine'
		new_event['connector_name'] = NAME
		new_event['event_type']     = 'check'
		new_event['source_type']    = 'component'
		new_event['component']      = INTERNAL_COMPONENT

		if 'resource' in new_event:
			del new_event['resource']

		self.count_alert(new_event, 1)

		logevent['source_type'] = 'resource'
		new_event['source_type'] = 'resource'

		for hostgroup in event.get('hostgroups', []):
			logevent['resource'] = hostgroup
			new_event['resource'] = hostgroup

			self.count_alert(new_event, 1)

			self.amqp.publish(logevent, cevent.get_routingkey(logevent), self.amqp.exchange_name_events)

	def count_sla(self, event, slatype, slaname, delay, value):
		now = int(time.time())

		def increment_SLA( event, slatype, slaname, delay, value, hostgroup=None ):

			meta_data = {'type': 'COUNTER', 'co': INTERNAL_COMPONENT}

			if hostgroup != None:
				meta_data['re'] = hostgroup
			#last_state_change field is updated in event store, so here we have no real previous date
			if 'previous_state_change_ts' in event:
				compare_date = event['previous_state_change_ts']
			else:
				compare_date = event['last_state_change']

			ack = self.entities.find_one({
				'type': 'ack',
				'timestamp': {
					'$gt': compare_date,
					'$lt': now
				}
			})

			#when time elapsed exceed deadline, produce an SLA out increment or nok if alert was ack else, auto increment if only event was not ack ack, SLA ok otherwise

			sla_states = {
				'out': 0,
				'nok': 0,
				'ok': 0
			}

			# set increment 1 depending on computation rules
			if delay < (now - compare_date):
				if ack:
					sla_states['nok'] = 1
				else:
					sla_states['out'] = 1
			elif ack:
				sla_states['ok'] = 1
			else:
				# spontaneous solve case
				if event['state'] == 0:
					meta_data_auto = {
						'type': 'COUNTER',
						'co': INTERNAL_COMPONENT,
						'me': "cps_sla_autosolve_{0}".format(slaname) ,
					}
					if hostgroup != None:
						meta_data_auto['re'] = hostgroup
					self.increment_counter(meta_data_auto, 1)

			# increment all counts with computed value

			for sla_state in sla_states:
				meta_data['me'] = 'cps_sla_{0}_{1}_{2}'.format(
					slatype,
					slaname,
					sla_state
				)
				self.increment_counter(meta_data, sla_states[sla_state])


		for hostgroup in event.get('hostgroups', []):
			increment_SLA( event, slatype, slaname, delay, value, hostgroup)

		increment_SLA( event, slatype, slaname, delay, value)

	def count_by_crits(self, event, value):

		if 'previous_state'in event and event['state'] == 0 and event.get('state_type', 1) == 1:

			sla_field = event.get(self.MACRO, None)

			if sla_field and sla_field in self.crits:
				if event['previous_state'] == 1:
					self.count_sla(event, 'warn', sla_field, self.crits[sla_field], value)

				if event['previous_state'] == 2:
					self.count_sla(event, 'crit', sla_field, self.crits[sla_field], value)

			# Update others
			meta_data = {'type': 'COUNTER', 'co': INTERNAL_COMPONENT }

			for _crit in self.crits:
				# Update warning counters
				if _crit != sla_field:
					for slatype in ['ok', 'nok', 'out']:
						meta_data['me'] = 'cps_sla_warn_{0}_{1}'.format(_crit, slatype)
						self.increment_counter(meta_data, 0)


	def count_alert(self, event, value):
		component = event['component']
		resource = event.get('resource', None)
		tags = event.get('tags', [])
		state = event['state']
		state_type = event.get('state_type', 1)

		# Update cps_statechange{,_0,_1,_2,_3} for component/resource

		meta_data = {
			'type': 'COUNTER',
			'co': component,
			'tg': tags
		}

		if resource:
			meta_data['re'] = resource

		meta_data['me'] = "cps_statechange"
		self.increment_counter(meta_data, value)

		meta_data['me'] = "cps_statechange_nok"
		cvalue = value if state != 0 else 0
		self.increment_counter(meta_data, cvalue)

		for cstate in [0, 1, 2, 3]:
			cvalue = value if cstate == state else 0

			meta_data['me'] = "cps_statechange_{0}".format(cstate)
			self.increment_counter(meta_data, cvalue)

		# Update cps_statechange_{hard,soft}

		for cstate_type in [0, 1]:
			cvalue = value if cstate_type == state_type and state != 0 else 0

			meta_data['me'] = "cps_statechange_{0}".format(
				'hard' if cstate_type == 1 else 'soft'
			)

			self.increment_counter(meta_data, cvalue)

	def count_by_type(self, event, value):
		state = event['state']
		def count_by_type_hostgroups(event, value, hostgroup=None):
			#Shortcut
			def increment(increment_type, value, hostgroup):
				metas = {
					'type': 'COUNTER',
					'co': INTERNAL_COMPONENT,
					'tg': event.get('tags', []),
					'me': "cps_statechange_{0}".format(increment_type)
				}
				if hostgroup:
					metas['re'] = hostgroup
				self.increment_counter(metas, value)

			#Keep only logic. increment component if on error
			if event['source_type'] == 'component':
				if state != 0:
					increment('component', value, hostgroup)
				else:
					increment('component', 0, hostgroup)

			# increment resource if in error. status depends on it s component. increment resource by component if in error by component
			if event['source_type'] == 'resource':

				component_problem = False
				if cevent.is_component_problem(event):
					component_problem = True
					increment('resource_by_component', value, hostgroup)
				else:
					increment('resource_by_component', 0, hostgroup)

				if state != 0 or component_problem:
					increment('resource', value, hostgroup)
				else:
					increment('resource', 0, hostgroup)

			meta_data = {
				'type': 'COUNTER',
				'co': INTERNAL_COMPONENT,
				'tg': event.get('tags', [])
			}
			if hostgroup:
				meta_data['re'] = hostgroup
			# Update cps_alerts_not_ack

			if state != 0 and event['state_type'] == 1:
				ackhost = cevent.is_host_acknowledged(event)
				cvalue0 = int(not ackhost)
				cvalue1 = int(not not ackhost)

				meta_data['me'] = 'cps_alerts_not_ack'
				self.increment_counter(meta_data, cvalue0)

				meta_data['me'] = 'cps_alerts_ack'
				self.increment_counter(meta_data, 0)

				meta_data['me'] = 'cps_alerts_ack_by_host'
				self.increment_counter(meta_data, cvalue1)

		for hostgroup in event.get('hostgroups', []):
			count_by_type_hostgroups(event, value, hostgroup)

		count_by_type_hostgroups(event, value)

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

		if event['rk'] in self.comments:
			self.increment_counter({
				'type': 'COUNTER',
				'co': INTERNAL_COMPONENT,
				'tg': event.get('tags', []),
				'me': 'cps_alerts_mass_ack'
			}, 1)
			# This event is exclided from counts because it's ack contained a special comment that matched withs configuration ones.
			return event

		if 'downtime' not in event or not event['downtime']:

			validation = event['event_type'] in self.listened_event_type
			validation = validation and event['component'] not in ['derogation', INTERNAL_COMPONENT]

			if validation:

					self.update_global_counter(event)
					self.count_by_crits(event, 1)

					# By name
					self.count_alert(event, 1)

					# By Type and ACK
					self.count_by_type(event, 1)

					# By tags (selector)
					self.count_by_tags(event, 1)

			return event
