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

import clogging, time
from cstorage import get_storage
from caccount import caccount
from crecord import crecord

from ctools import legend
from ctools import uniq
from ctools import parse_perfdata

legend_type = ['soft', 'hard']

class carchiver(object):
	def __init__(self, namespace, storage=None, autolog=False):

		self.logger = clogging.getLogger('carchiver-{0}'.format(namespace))

		self.namespace = namespace
		self.namespace_log = namespace + '_log'

		self.autolog = autolog

		self.logger.debug("Init carchiver on %s" % namespace)

		self.account = caccount(user="root", group="root")

		if not storage:
			self.logger.debug(" + Get storage")
			self.storage = get_storage(namespace=namespace)
		else:
			self.storage = storage

		self.collection = self.storage.get_backend(namespace)

	def check_event(self, _id, event):
		changed = False
		new_event = False

		self.logger.debug(" + Event:")

		state = event['state']
		state_type = event['state_type']

		self.logger.debug("   - State:\t\t'%s'" % legend[state])
		self.logger.debug("   - State type:\t'%s'" % legend_type[state_type])

		now = int(time.time())
		
		event['timestamp'] = event.get('timestamp', now)
		
		try:
			# Get old record
			#record = self.storage.get(_id, account=self.account)

			devent = self.collection.find_one(_id, fields={'state': 1, 'state_type': 1, 'last_state_change': 1, 'perf_data_array': 1})
			
			self.logger.debug(" + Check with old record:")
			old_state = devent['state']
			old_state_type = devent['state_type']
			
			event['last_state_change'] = devent.get('last_state_change', event['timestamp'])

			self.logger.debug("   - State:\t\t'%s'" % legend[old_state])
			self.logger.debug("   - State type:\t'%s'" % legend_type[old_state_type])

			if state != old_state:
				event['previous_state'] = old_state

			if state != old_state or state_type != old_state_type:
				self.logger.debug(" + Event has changed !")
				changed = True
			else:
				self.logger.debug(" + No change.")
			
			try:
				event = self.merge_perf_data(devent, event)
			except Exception, err:
				self.logger.warning("merge_perf_data: %s" % err)

		except:
			# No old record
			self.logger.debug(" + New event")
			changed = True
			new_event = True
		
		if changed:
			event['last_state_change'] = event.get('timestamp', now)
		
		if new_event:
			self.store_new_event(_id, event)
		else:
			self.store_update_event(_id, event)

		mid = None
		if changed and self.autolog:
			mid = self.log_event(_id, event)			

		return mid
		
	def merge_perf_data(self, old_event, new_event):
		old_event['perf_data_array'] = old_event.get('perf_data_array', [])
		new_event['perf_data_array'] = new_event.get('perf_data_array', [])
		
		if new_event['perf_data_array'] != []:
			perf_data_array = old_event['perf_data_array']
			
			new_metrics = [ perf['metric'] for perf in new_event['perf_data_array'] ]
			old_metrics = [ perf['metric'] for perf in old_event['perf_data_array'] ]
						
			if new_metrics == old_metrics:
				new_event['perf_data_metrics'] = new_metrics
				return new_event
			
			new_event['perf_data_metrics'] = uniq(new_metrics + old_metrics)
			
			for new_metric in new_metrics:
				if new_metric in old_metrics:
					perf_data_array[old_metrics.index(new_metric)] = new_event['perf_data_array'][new_metrics.index(new_metric)]
				else:
					perf_data_array.append(new_event['perf_data_array'][new_metrics.index(new_metric)])
						
			new_event['perf_data_array'] = perf_data_array
			
		
		
		return new_event

	def store_new_event(self, _id, event):
		record = crecord(event)
		record.type = "event"
		record.chmod("o+r")
		record._id = _id

		self.storage.put(record, namespace=self.namespace, account=self.account)

	def store_update_event(self, _id, event):
		self.collection.update({'_id': _id}, {"$set": event}, safe=True)
	
	def log_event(self, _id, event):
		self.logger.debug("Log event '%s' in %s ..." % (_id, self.namespace_log))
		record = crecord(event)
		record.type = "event"
		record.chmod("o+r")
		record.data['event_id'] = _id
		record._id = _id + '.' + str(time.time())

		self.storage.put(record, namespace=self.namespace_log, account=self.account)
		return record._id

	def get_logs(self, _id, start=None, stop=None):
		return self.storage.find({'event_id': _id}, namespace=self.namespace_log, account=self.account)

	def remove_all(self):
		self.logger.debug("Remove all logs and state archives")

		self.storage.drop_namespace(self.namespace)
		self.storage.drop_namespace(self.namespace_log)
		
	
