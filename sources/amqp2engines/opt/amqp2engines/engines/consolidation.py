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

#states_str = ("Ok", "Warning", "Critical", "Unknown", "Undetermined")
#states = {0: 0, 1:0, 2:0, 3:0, 4:0}

states_str = ("Ok", "Warning", "Critical", "Unknown")
states = {0: 0, 1:0, 2:0, 3:0}

class engine(cengine):
	def __init__(self, *args, **kargs):
		self.metrics_list = {}
		self.timestamp = { } 
		self.manager = pyperfstore2.manager(logging_level=logging.INFO)
		self.beat_interval = 300 
	
		#for debug
		#self.beat_interval = 5 
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.default_interval = 300
		self.records = { } 
		
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.manager = pyperfstore2.manager(logging_level=self.logging_level)
		self.load_consolidation()
		self.beat()
	def beat(self):
		non_loaded_records = self.storage.find({ '$and' : [{ 'crecord_type': 'consolidation' }, {'loaded': { '$ne' : 'true'} } ] }, namespace="object" )
		
		if len(non_loaded_records) > 0  :
			for i in non_loaded_records :
				self.load(i)
		for _id in self.records.keys() :
			exists = self.storage.find({ '_id': _id } )
			if len(exists) == 0  :
				del(self.records[_id])
			elif len(exists) == 1:
				rec = exists[0].dump()
				self.records[_id]['enable'] = rec.get('enable')

		for record in self.records.values():
			interval = record.get('interval', self.default_interval)
			if  int(interval) < ( int(time.time()) - self.timestamp[record.get('_id')]) and ( record.get('enable') == "true" or record.get('enable') == True ) :
				tfilter = json.loads(record.get('mfilter'))
				metric_list = self.manager.store.find(mfilter=tfilter)
				values = []
				list_fn = record.get('type', False)
				if isinstance(list_fn, str) or isinstance(list_fn, unicode) :
					list_fn = [ list_fn ] 
				for metric in metric_list :
					m = metric.get('d')
					if ( len(m) >0 ) :
						values.append( m[-2:-1] ) 
				
				if list_fn and len(values) > 0 :
					list_perf_data = []
					for i in list_fn :
						if i == 'mean':
							fn = lambda x: sum(x) / len(x)
						elif i == 'min' :
							fn = lambda x: min(x)
						elif i == 'max' :
							fn = lambda x: max(x)
						elif i == 'sum':
							fn = lambda x: sum(x)
						elif i == 'delta':
							fn = lambda x: x[0] - x[-1]
						resultat = list()
						try :
							resultat = pyperfstore2.utils.aggregate_series(values, fn)
						except NameError:
							self.logger.info('la fonction '+i+' est inexistante')
							self.storage.update(record.get('_id'), {'output_engine': "function "+i+" does not exists"  } )
						if len(resultat) > 0 :
							list_perf_data.append({ 'metric' : i, 'value' : resultat[0][1], "unit": None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' } ) 
							event = cevent.forger(
								connector ="consolidation",
								connector_name = "engine",
								event_type = "consolidation",
								source_type = "resource",
								component = record['component'],
								resource=record['resource'],
								state=0,
								state_type=0,
								output="",
								long_output="",
		        			                perf_data=None,
			                		        perf_data_array=list_perf_data,
		        			                display_name=record['crecord_name'][0]
							)	
							rk = cevent.get_routingkey(event)
							self.amqp.publish(event, rk, self.amqp.exchange_name_events)
							self.storage.update(record.get('_id'), {'output_engine': datetime.now().strftime('%Y-%m-%d %H:%M:%S')+" : Computation done. Next Computation in "+str(interval)+" s"  } )
						else:
							self.storage.update(record.get('_id'), {'output_engine': "No result"  } )
				else:
					self.storage.update(record.get('_id'), {'output_engine': "No input values"  } )
				self.timestamp[record.get('_id')] = int(time.time())
		
		
	def load (self, rec ) :
		self.logger.debug('load')
		record = rec.dump()
		rec.loaded = True
		self.storage.update(record.get('_id'), {'loaded': 'true' })
		if record.get('mfilter', False) :
			self.timestamp[record.get('_id')] = int(time.time())
			tfilter = json.loads(record.get('mfilter'))
			metric_list = self.manager.store.find(mfilter=tfilter )
			nb_items = metric_list.count()
			self.storage.update(record.get('_id'), {'nb_items': nb_items } )
			self.storage.update(record.get('_id'), {'output_engine': "Correctly Load"  } )
			event = cevent.forger(
					connector = "consolidation",
					connector_name = "engine",
					event_type = "check",
					source_type="resource",
					component=record['crecord_name'][1],
					resource=record['crecord_name'][2],
					state=0,
					state_type=1,
					output="",
					long_output="",
		                        perf_data=None,
                		        perf_data_array=None,
		                        display_name=record['crecord_name'][0]
			)
			rk = cevent.get_routingkey(event)
			self.records[record.get('_id')] = record
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		else:
			self.storage.update(record.get('_id'), {'output_engine': "Impossible to load : no filter defined"  } )

	def load_consolidation(self) :
		records = self.storage.find({ '$and' :[ {'crecord_type': 'consolidation'}] }, namespace="object")
		for i in records :
			self.load(i)
				
			
			
			
	def unload_consolidation(self):
		record_list = self.storage.find({ '$and': [{'crecord_type': 'consolidation' }, {'loaded':'true'}]}, namespace="object")
		for i in record_list :
			self.storage.update(i._id, {'loaded': 'false' })
			self.storage.update(i._id, {'output_engine': "Correctly Unload"  } )
