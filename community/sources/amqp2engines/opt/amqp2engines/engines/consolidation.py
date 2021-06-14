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
import pyperfstore2.utils
import cevent
import logging
import json

import time
from datetime import datetime
from ctools import internal_metrics, roundSignifiantDigit


class engine(cengine):
	etype = 'consolidation'

	def __init__(self, *args, **kargs):
		super(engine, self).__init__(*args, **kargs)

		self.metrics_list = {}
		self.timestamps = {} 
		self.default_interval = 60

		self.thd_warn_sec_per_evt = 8
		self.thd_crit_sec_per_evt = 10


	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.manager = pyperfstore2.manager(logging_level=self.logging_level)
				
		self.beat()

	def beat(self):
		self.logger.debug('Consolidation BEAT')

	def consume_dispatcher(self,  event, *args, **kargs):
		self.logger.debug("Consolidate metrics:")

		now = time.time()
		beat_elapsed = 0

		record = self.get_ready_record(event)
		if record:	
			record = record.dump()		

			_id = record.get('_id')
			name = record.get('crecord_name')

			aggregation_interval = record.get('aggregation_interval')

			self.logger.debug("'%s':" % name)
			self.logger.debug(" + interval: %s" % aggregation_interval)

			last_run = record.get('consolidation_ts', now)

			elapsed = now - last_run

			self.logger.debug(" + elapsed: %s" % elapsed)

			mfilter = record.get('mfilter')

			#Nothing to do if no filter set
			if mfilter and (elapsed == 0 or elapsed >= aggregation_interval):
				self.logger.debug("Step 1: Select metrics")

				mfilter = json.loads(mfilter)

				self.logger.debug(' + mfilter: %s' % mfilter)

				and_clause = [mfilter, {'me': {'$nin':internal_metrics}}]

				# Exclude internal metrics
				mfilter = {'$and': and_clause}

				metric_list = self.manager.store.find(mfilter=mfilter)

				self.logger.debug(" + %s metrics found" % metric_list.count())

				if not metric_list.count():
					self.storage.update(_id, { 'output_engine': "No metrics, check your filter" })
				else:

					aggregation_method = record.get('aggregation_method')
					self.logger.debug(" + aggregation_method: %s" % aggregation_method)

					consolidation_methods = record.get('consolidation_method')
					if not isinstance(consolidation_methods, list):
						consolidation_methods = [ consolidation_methods ]

					self.logger.debug(" + consolidation_methods: %s" % consolidation_methods)

					mType = mUnit = mMin = mMax = None
					sum_in_consolidation_methods = 'sum' in consolidation_methods
					maxSum = 0
					# Get metrics
					metrics = []
					for index, metric in enumerate(metric_list):
						if index == 0 :
							#mType = metric.get('t')
							mMin = metric.get('mi')
							mMax = metric.get('ma')
							mUnit = metric.get('u')
							if sum_in_consolidation_methods and mMax is not None:
								maxSum = mMax
						else:
							mi = metric.get('mi')
							if mi is not None and (mMin is None or mi < mMin):
								mMin = mi
							ma = metric.get('ma')
							if ma is not None:
								if mMax is None or ma > mMax:
									mMax = ma
								if sum_in_consolidation_methods and mMax is not None:
									maxSum += ma
							if metric.get('u') != mUnit :
								self.logger.warning("%s: too many units" % name)
								output_message = "warning : too many units"

						self.logger.debug(' + %s , %s , %s, %s' % (
							metric.get('_id'),
							metric.get('co'),
							metric.get('re',''),
							metric.get('me'))
						)

						metrics.append(metric.get('_id'))

					self.logger.debug(' + mMin: %s' % mMin)
					self.logger.debug(' + mMax: %s' % mMax)
					self.logger.debug(' + mUnit: %s' % mUnit)

					self.logger.debug("Step 2: Aggregate (%s)" % aggregation_method)

					# Set time range
					tstart = last_run

					if elapsed == 0 or last_run < (now - 2 * aggregation_interval):
						tstart = now - aggregation_interval

					self.logger.debug(
						" + From: %s To %s "%
						(datetime.fromtimestamp(tstart).strftime('%Y-%m-%d %H:%M:%S'),
						datetime.fromtimestamp(time.time()).strftime('%Y-%m-%d %H:%M:%S'))
					)

					values = []
					for mid in metrics:
						points = self.manager.get_points(tstart=tstart, tstop=now, _id=mid)
						fn = self.get_math_function(aggregation_method)

						pValues = [point[1] for point in points]

						if not len(pValues):
							continue

						values.append(fn(pValues))

					self.logger.debug(" + %s values" % len(values))

					if not len(values):
						self.storage.update(_id, { 'output_engine': "No values, check your interval" })
					else:
						self.logger.debug("Step 3: Consolidate (%s)" % consolidation_methods)

						perf_data_array = []

						for consolidation_method in consolidation_methods:
							fn = self.get_math_function(consolidation_method)
							value = fn(values)

							self.logger.debug(" + %s: %s %s" % (consolidation_method, value, mUnit))

							perf_data_array.append({
								'metric' : consolidation_method,
								'value' : roundSignifiantDigit(value,3),
								"unit": mUnit,
								'max': maxSum if consolidation_method == 'sum' else mMax,
								'min': mMin,
								'type': 'GAUGE'
							})

						self.logger.debug("Step 4: Send event")

						event = cevent.forger(
							connector ="consolidation",
							connector_name = "engine",
							event_type = "consolidation",
							source_type = "resource",
							component = record['component'],
							resource=record['resource'],
							state=0,
							timestamp=now,
							state_type=1,
							output="Consolidation: '%s' successfully computed" % name,
							long_output="",
							perf_data=None,
							perf_data_array=perf_data_array,
							display_name=name
						)
						rk = cevent.get_routingkey(event)
						self.counter_event += 1
						self.amqp.publish(event, rk, self.amqp.exchange_name_events)

						self.timestamps[_id] = now

						self.logger.debug("Step 5: Update configuration")

						beat_elapsed = time.time() - now

						self.storage.update(_id, {
							'consolidation_ts': int(now),
							'nb_items': len(metrics),
							'output_engine': "Computation done in %.2fs (%s/%s)" % (beat_elapsed, len(values), len(metrics))
						})
			else:
				self.logger.debug("Not the moment to process this consolidation")

		else:
			self.logger.warning('Dispatch error: engine unable to load consolidation crecord properly')

		#set record free for dispatcher engine
		self.crecord_task_complete(_id)

		if not beat_elapsed:
			beat_elapsed = time.time() - now

		self.counter_worktime += beat_elapsed

	def get_math_function(self, name):
		if name == 'average' or name == 'mean':
			return lambda x: sum(x) / len(x)
		elif name == 'min' :
			return lambda x: min(x)
		elif name == 'max' :
			return lambda x: max(x)
		elif name == 'sum':
			return lambda x: sum(x)
		elif name == 'delta':
			return lambda x: x[0] - x[-1]
		elif name == 'last':
			return lambda x: x[len(x)-1]
		else:
			return None

	def post_run(self):
		pass
