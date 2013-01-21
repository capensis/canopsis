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

NAME="consolidation"

class engine(cengine):
	def __init__(self, *args, **kargs):
		self.metrics_list = {}
		self.timestamps = { } 
		self.manager = pyperfstore2.manager(logging_level=logging.INFO)
		self.beat_interval = 60
	
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.default_interval = 60
		self.records = { } 
		
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
			consolidation_interval = int(record.get('interval', self.default_interval))
			current_interval = int(time.time() - self.timestamps[record.get('_id')])

			self.logger.debug('current interval: %s , consolidation interval: %s' % (current_interval,consolidation_interval))
			if  current_interval >= consolidation_interval:
				self.logger.debug('Compute new consolidation for: %s' % record.get('crecord_name','Unkown'))

				output_message = None
				tfilter = json.loads(record.get('mfilter'))
				metric_list = self.manager.store.find(mfilter=tfilter)
				
				list_fn = record.get('type', False)

				if not isinstance(list_fn, list):
					list_fn = [ list_fn ] 

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

					self.logger.debug('Get last point for: %s' % metric.get('_id'))
					last_point = self.manager.get_last_point(_id=metric.get('_id'))
					if last_point:
						values.append([last_point])

				#self.logger.debug("type: %s" % type(values))
				#self.logger.debug(values)
				
				if list_fn and len(values) > 0 :
					list_perf_data = []
					for function_name in list_fn :
						fn = self.get_math_function(function_name)
						resultat = []
						try :
							resultat = pyperfstore2.utils.aggregate_series(values, fn)
						except NameError:
							self.logger.info('Function [%s] does not exist' % function_name)
							output_message = "warning : function [%s] does not exists" % function_name
						if len(resultat) > 0 :
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
								output="Consolidation: '%s' successfully loaded" % record['crecord_name'],
								long_output="",
								perf_data=None,
								perf_data_array=list_perf_data,
								display_name=record['crecord_name']
							)	
							rk = cevent.get_routingkey(event)
							self.amqp.publish(event, rk, self.amqp.exchange_name_events)

							if not output_message:
								engine_output = '%s : Computation done. Next Computation in %s s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S'),str(consolidation_interval))
								self.storage.update(record.get('_id'),{'output_engine':engine_output} )
							else:
								engine_output = '%s : Computation done but there are issues : "%s" . Next Computation in %s s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S'),output_message,str(consolidation_interval))
								self.storage.update(record.get('_id'), {'output_engine': engine_output} )
						else:
							if not output_message:
								self.storage.update(record.get('_id'), {'output_engine': "No result"  } )
							else:
								self.storage.update(record.get('_id'), {'output_engine': "there are issues : %s warning : No result" % output_message } )
				else:
					self.storage.update(record.get('_id'), {'output_engine': "No input values"  } )
				self.timestamps[record.get('_id')] = time.time()
		
		
	def load (self, rec ) :
		record = rec.dump()
		rec.loaded = True
		self.storage.update(record.get('_id'), {'loaded': True })
		if record.get('mfilter', False) :
			self.timestamps[record.get('_id')] = time.time()
			tfilter = json.loads(record.get('mfilter'))
			metric_list = self.manager.store.find(mfilter=tfilter )
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
					output="",
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