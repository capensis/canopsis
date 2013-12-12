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
#from pyperfstore import node
#from pyperfstore import mongostore
import pyperfstore2
import cevent
import logging

import time
from datetime import datetime

NAME="sla"

#states_str = ("Ok", "Warning", "Critical", "Unknown", "Undetermined")
#states = {0: 0, 1:0, 2:0, 3:0, 4:0}

states_str = ("Ok", "Warning", "Critical", "Unknown")
states = {0: 0, 1:0, 2:0, 3:0}

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		
		self.create_queue = False
				
		self.beat_interval =  900
		
		# For debug
		#self.beat_interval =  60
		
		self.resource = "sla"
		
		self.thd_warn_sla_timewindow = 98
		self.thd_crit_sla_timewindow = 95
		self.default_sla_timewindow = 60*60*24 # 1 day
		
		self.default_sla_output_tpl="{cps_pct_by_state_0}% Ok, {cps_pct_by_state_1}% Warning, {cps_pct_by_state_2}% Critical, {cps_pct_by_state_3}% Unknown"

	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.manager = pyperfstore2.manager(logging_level=self.logging_level)
		self.beat()
		
	def split_state(self, value):
		# cps_state = state * 100 + state_type * 10 + state_extra
		try:
			value = str(value)
			
			if len(value) == 2:
				value = "0%s" % value
				
			state = int(value[0])
			state_type = int(value[1])
			extra = int(value[2])
		except Exception, err:
			self.logger.error("Invalid value format: %s (%s)" % (value, err))
			raise Exception("Invalid value format: %s (%s)" % (value, err))
			
		return (state, state_type, extra)
	
	def get_states(self, name, metric, start, stop):
		metric_name = '%s%s' % (name,metric)
		try:
			points = self.manager.get_points(name=metric_name, tstart=start, tstop=stop)
		except:
			return []

		return points
		
	def get_rk(self, name):
		return "sla.engine.sla.resource.%s.sla" % name
		
	def get_name(self,name):
		return '%ssla' % name
	
	
	def calcule_sla(self, _id, config):
		
		self.logger.debug(" + Calcul SLA of '%s'" % config['name'])
		
		sla_timewindow = config.get('sla_timewindow', self.default_sla_timewindow)
		thd_warn_sla_timewindow = config.get('thd_warn_sla_timewindow', self.thd_warn_sla_timewindow)
		thd_crit_sla_timewindow = config.get('thd_crit_sla_timewindow', self.thd_crit_sla_timewindow)
		sla_output_tpl = config.get('sla_output_tpl', self.default_sla_output_tpl)
		sla_timewindow_doUnknown = config.get('sla_timewindow_doUnknown', True)
	
		# Prevent empty string
		if not isinstance(sla_timewindow, int):
			self.logger.warning("%s: Invalid 'sla_timewindow': %s" % (_id, sla_timewindow))
			sla_timewindow = self.default_sla_timewindow
		
		if not thd_warn_sla_timewindow:
			thd_warn_sla_timewindow = self.thd_warn_sla_timewindow
		if not thd_crit_sla_timewindow:
			thd_crit_sla_timewindow = self.thd_crit_sla_timewindow
		if sla_output_tpl == "":
			sla_output_tpl = self.default_sla_output_tpl	
			
		thd_warn_sla_timewindow = float(thd_warn_sla_timewindow)
		thd_crit_sla_timewindow = float(thd_crit_sla_timewindow)
		
		# For debug
		#sla_timewindow = 60*60

		self.logger.debug("   + sla timewindow:    %s" % sla_timewindow)
	
		stop = int(time.time())
		start = stop - sla_timewindow
		
	
		self.logger.debug("   + sla doUnknown:     %s" % sla_timewindow_doUnknown)
		self.logger.debug("   + start:             %s (%s)" % (datetime.utcfromtimestamp(start), start))
		self.logger.debug("   + stop:              %s (%s)" % (datetime.utcfromtimestamp(stop), stop))
		
		try:
			points = self.manager.get_points(name="%s%s" % (config['name'], 'cps_state'), tstart=start, tstop=stop, add_prev_point=True, add_next_point=True)
		except Exception, err:
			self.logger.error("Error when 'get_points': %s" % err)
			points = []
		
		self.logger.debug("   + Nb points:         %s" % len(points))
		
		if len(points) < 2:
			self.logger.debug("     + Need more points")
			return
		
		# For Debug
		#points.insert(0, [start-75, 210])
		#points.append([stop+75, 210])

		first_point = points.pop(0)
		last_point =  points.pop(len(points)-1)
		
		self.logger.debug("   + First point:       %s" % datetime.utcfromtimestamp(first_point[0]))
		self.logger.debug("   + Last point:        %s" % datetime.utcfromtimestamp(last_point[0]))
		self.logger.debug("   + Total:             %s" % (last_point[0] - first_point[0]))
		
		start_undetermined_time =  first_point[0] - start
		stop_undetermined_time = stop - last_point[0]
		
		self.logger.debug("   + Start utime:       %s" % start_undetermined_time)
		if start_undetermined_time < 0:
			self.logger.debug("     + Sample start too soon, adjust timestamp")
			first_point[0] = start
			
		self.logger.debug("   + Stop utime:        %s" % stop_undetermined_time)
		if start_undetermined_time < 0:
			self.logger.debug("     + Sample finish too late, adjust timestamp")
			last_point[0] = stop	
		
		self.logger.debug("   + Total utime:       %s" % (stop_undetermined_time + start_undetermined_time))
		
		states_sum = states.copy()
		total = 0
		
		if sla_timewindow_doUnknown and start_undetermined_time > 0:
			self.logger.debug("     + Set %s seconds to Unknown state" % start_undetermined_time)
			states_sum[3] = start_undetermined_time
			total += start_undetermined_time
		
		(state, state_type, extra) = self.split_state(first_point[1])
		timestamp = first_point[0]
		self.logger.debug("   + Initial State:     %s" % state)
		
		# Sum all duration by state
		for point in points:
			duration = point[0] - timestamp	
			states_sum[state] += duration
			total += duration
			
			(state, state_type, extra) = self.split_state(point[1])
			timestamp = point[0]
			
		# Finish by last point
		duration = last_point[0] - timestamp
		states_sum[state] += duration
		total += duration		
			
		self.logger.debug("   + States sum:        %s" % states_sum)
		self.logger.debug("   + Total:             %s" % total)
		
		# for event
		output_data = {}
		perf_data_array = []
		
		states_pct = states.copy()
		for state in states_sum:
			if states_sum[state] != 0:
				states_pct[state] = round((states_sum[state] * 100)/float(total), 3)
			else:
				states_pct[state] = 0
				
			metric = 'cps_pct_by_state_%s' % state
			output_data[metric] = states_pct[state]
			perf_data_array.append({"metric": metric, "value": states_pct[state], "max": 100, "unit": "%"})
				
		self.logger.debug("   + States pct:        %s" % states_pct)
		
		self.logger.debug("   + Event:")
		
		# Calcul state
		state = 0
		if states_pct[0] < thd_warn_sla_timewindow:
			state = 1
		if states_pct[0] < thd_crit_sla_timewindow:
			state = 2
		
		self.logger.debug("     + State:     %s (%s)" % (states_str[state], state))
		
		## Build Event
		
		# Fill output (for event)
		output = sla_output_tpl
		if output_data:
			for key in output_data:
				output = output.replace("{%s}" % key, str(output_data[key]))
				
		
		self.logger.debug("     + Output:    %s" % output)
		self.logger.debug("     + Perfdata:  %s" % perf_data_array)
		
		# Send AMQP Event
		event = cevent.forger(
			connector = "sla",
			connector_name = "engine",
			event_type = "sla",
			source_type="resource",
			component=config['name'],
			resource="sla",
			state=state,
			state_type=1,
			output=output,
			long_output="",
			perf_data=None,
			perf_data_array=perf_data_array,
			display_name=config.get('display_name', None)
		)
		
		# Extra fields
		event['selector_id'] = config['_id']
		event['selector_rk'] = config['rk']
		
		rk = self.get_rk(config['name'])
		
		self.logger.debug("Publish event on %s" % rk)
		self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		
		# Update selector record
		self.storage.update(_id, {'sla_timewindow_lastcalcul': stop, 'sla_timewindow_perfdata': perf_data_array, 'sla_state': event['state'], 'sla_rk': rk})
	
	def beat(self):
		self.logger.debug('BEAT sla')	
	
		
	def consume_dispatcher(self,  event, *args, **kargs):
	
		self.logger.debug('entered in sla consume dispatcher')
		# Gets crecord from amqp distribution
		record = self.get_ready_record(event)
		
		if record:	
			_id = event['_id']

			record_data 		= record.data
			record_data['name'] = record.name
			record_data['_id'] 	= record._id

			self.logger.debug("Load selector '%s' (%s)" % (record_data['name'], _id))
		
			sla_id = self.calcule_sla(_id, record_data)
		
			self.counter_event += 1
			self.crecord_task_complete(event['_id'])

				
