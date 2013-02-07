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
from ctools import internal_metrics

NAME="consolidation"

class engine(cengine):
	def __init__(self, *args, **kargs):
		self.metrics_list = {}
		self.timestamps = {} 
		self.records = {} 
		self.default_interval = 60

		self.manager = pyperfstore2.manager(logging_level=logging.INFO)	
		cengine.__init__(self, name=NAME, *args, **kargs)

	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.manager = pyperfstore2.manager(logging_level=self.logging_level)
		self.load_consolidation()
		self.beat()

	def beat(self):
		non_loaded_records = self.storage.find({ '$and' : [{ 'crecord_type': 'consolidation' },{'enable': True}, {'loaded': { '$ne' : True} } ] }, namespace="object" )

		if len(non_loaded_records) > 0  :
			for item in non_loaded_records :
				self.logger.info("New consolidation found '%s', load" % item.name)
				self.load(item)

		for _id in self.records.keys() :
			exists = self.storage.find_one({ '_id': _id } )
			if not exists:
				self.logger.info("%s deleted, remove from record list" % self.records[_id]['crecord_name'])
				del(self.records[_id])

		for record in self.records.values():
			consolidation_last_timestamp = self.timestamps[record.get('_id')]

			aggregation_interval = record.get('aggregation_interval', self.default_interval)
			current_interval = int(time.time()) - consolidation_last_timestamp

			self.logger.debug('current interval: %s , consolidation interval: %s' % (current_interval,aggregation_interval))
			if  current_interval >= aggregation_interval:
				self.logger.debug('Compute new consolidation for: %s' % record.get('crecord_name','No name found'))

				output_message = None
				mfilter = json.loads(record.get('mfilter'))
				mfilter = {'$and': [mfilter, {'me': {'$nin':internal_metrics}}]}
				self.logger.debug('the mongo filter is: %s' % mfilter)
				metric_list = self.manager.store.find(mfilter=mfilter)
				self.logger.debug('length of matching metric list is: %i' % metric_list.count())

				
				first_aggr_function = record.get('first_aggregation_type', False)
				second_aggr_function = record.get('second_aggregation_type', False)

				if not isinstance(second_aggr_function, list):
					second_aggr_function = [ second_aggr_function ] 

				mType = mUnit = mMin = mMax = None
				values = []

				for index,metric in enumerate(metric_list) :
					if  index == 0 :
						mType = metric.get('t')
						mMin = metric.get('mi')
						mMax = metric.get('ma')
						mUnit = metric.get('u')
						#mCrit = metric.get('tc')
						#mWarn = metric.get('tw')
					else:
						if  metric.get('mi') < mMin :
							mMin = metric.get('mi')
						if metric.get('ma') > mMax :
							mMax = metric.get('ma')
						if metric.get('u') != mUnit :
							output_message = "warning : too many units"
						if  mType != metric.get('t') :
							output_message = "warning : too many metrics type"

					self.logger.debug(' + Get points for: %s , %s , %s, %s' % (metric.get('_id'),metric.get('co'),metric.get('re',''),metric.get('me')))

					if int(time.time()) - aggregation_interval <= consolidation_last_timestamp + 60:
						tstart = consolidation_last_timestamp
						self.logger.debug('   +   Use original tstart: %i' % consolidation_last_timestamp)
					else:
						tstart = int(time.time()) - aggregation_interval
						self.logger.debug('   +   new tstart: %i' % tstart)

					list_points = self.manager.get_points(tstart=tstart, _id=metric.get('_id'))
					self.logger.debug('   +   Values on interval:')
					self.logger.debug(list_points)

					if list_points:
						fn = self.get_math_function(first_aggr_function)
						point_timestamp = int(time.time()) - current_interval/2
						if fn:
							point_value = fn([value[1] for value in list_points])
						else:
							point_value = list_points[len(list_points)-1][1]
						values.append([[point_timestamp,point_value]])
				
				self.logger.debug('   +   Summary of horizontal aggregation:')
				self.logger.debug(values)

				if not second_aggr_function:
					self.storage.update(record.get('_id'), {'output_engine': "No second aggregation function given"  } )
					return

				if len(values) == 0 :
					self.logger.debug('  +  No values')
					self.storage.update(record.get('_id'), {
															'output_engine': "No input values",
															'consolidation_ts':int(time.time())
															})
					self.timestamps[record.get('_id')] = int(time.time())
					return

				list_perf_data = []
				for function_name in second_aggr_function :
					fn = self.get_math_function(function_name)

					if not fn:
						self.logger.debug('No function given for second aggregation')
						self.storage.update(record.get('_id'), {'output_engine': "No function given for second aggregation"})
						return

					resultat = pyperfstore2.utils.consolidation(values, fn)
	
					if len(resultat) == 0 :
						if not output_message:
							self.storage.update(record.get('_id'), {'output_engine': "No result"  } )
						else:
							self.storage.update(record.get('_id'), {'output_engine': "there are issues : %s warning : No result" % output_message } )

					self.logger.debug(' + Result of aggregation for "%s": %f' % (function_name,resultat[0][1]))

					list_perf_data.append({ 'metric' : function_name, 'value' : resultat[0][1], "unit": mUnit, 'max': mMax, 'min': mMin, 'warn': None, 'crit': None, 'type': mType } ) 
					event = cevent.forger(
						connector ="consolidation",
						connector_name = "engine",
						event_type = "consolidation",
						source_type = "resource",
						component = record['component'],
						resource=record['resource'],
						state=0,
						timestamp=resultat[0][0],
						state_type=0,
						output="Consolidation: '%s' successfully computed" % record.get('crecord_name','No name'),
						long_output="",
						perf_data=None,
						perf_data_array=list_perf_data,
						display_name=record['crecord_name']
					)	
					rk = cevent.get_routingkey(event)
					self.amqp.publish(event, rk, self.amqp.exchange_name_events)

					self.logger.debug('The following event was sent:')
					self.logger.debug(event)

					if not output_message:
						engine_output = '%s : Computation done. Next Computation in %s s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S'),str(aggregation_interval))
						self.storage.update(record.get('_id'),{'output_engine':engine_output} )
					else:
						engine_output = '%s : Computation done but there are issues : "%s" . Next Computation in %s s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S'),output_message,str(aggregation_interval))
						self.storage.update(record.get('_id'), {'output_engine': engine_output} )

				self.storage.update(record.get('_id'), {'consolidation_ts':int(time.time())})
				self.timestamps[record.get('_id')] = int(time.time())
		
		
	def load (self, rec ) :
		record = rec.dump()
		rec.loaded = True
		self.storage.update(record.get('_id'), {'loaded': True })
		if record.get('mfilter', False) :
			self.timestamps[record.get('_id')] = record.get('consolidation_ts',int(time.time()))
			mfilter = json.loads(record.get('mfilter'))
			metric_list = self.manager.store.find(mfilter=mfilter )
			nb_items = metric_list.count()
			self.storage.update(record.get('_id'), {
													'output_engine': "Correctly Loaded",
													'nb_items': nb_items
													} )
			event = cevent.forger(
					connector = "consolidation",
					connector_name = "engine",
					event_type = "check",
					source_type="resource",
					component=record['component'],
					resource=record['resource'],
					state=0,
					state_type=1,
					output="Consolidation %s successfully loaded" % record.get('crecord_name','No name'),
					long_output="",
					perf_data=None,
					perf_data_array=None,
					display_name=record['crecord_name']
			)
			rk = cevent.get_routingkey(event)
			self.records[record.get('_id')] = record
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		else:
			self.storage.update(record.get('_id'), {'output_engine': "Impossible to load : no filter defined"  } )

	def load_consolidation(self) :
		records = self.storage.find({ '$and' :[ {'crecord_type': 'consolidation'},{'enable': True}] }, namespace="object")
		for item in records :
			self.load(item)

		self.logger.info('Load %i consolidations' % len(records))
				
	def unload_consolidation(self):
		records = self.storage.find({ '$and': [{'crecord_type': 'consolidation' }, {'loaded':True}]}, namespace="object")
		for item in records :
			self.storage.update(item._id, {
										'output_engine': "Correctly Unload",
										'loaded': False
										} )

		self.logger.info('Unload %i consolidations' % len(records))

	def get_math_function(self, name):
		if name == 'mean':
			return lambda x: sum(x) / len(x)
		elif name == 'min' :
			return lambda x: min(x)
		elif name == 'max' :
			return lambda x: max(x)
		elif name == 'sum':
			return lambda x: sum(x)
		elif name == 'delta':
			return lambda x: x[0] - x[-1]
		else:
			return None
