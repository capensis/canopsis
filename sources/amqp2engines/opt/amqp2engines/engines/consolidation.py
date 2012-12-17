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
		print "init"
		self.metrics_list = {}
		self.timestamp = { } 
		self.manager = pyperfstore2.manager(logging_level=logging.INFO)
		self.beat_interval = 5 
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.default_interval = 10
		
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.manager = pyperfstore2.manager(logging_level=self.logging_level)
		self.load_consolidation()
		self.beat()
	def beat(self):
		self.logger.debug('beat')
		non_loaded_records = self.storage.find({ '$and' : [{ 'crecord_type': 'consolidation' }, {'loaded': { '$ne' : 'true'} } ] }, namespace="object" )
		
		if (len(non_loaded_records) > 0 ) :
			for i in non_loaded_records :
				self.load(i)
		self.logger.debug(self.records)
		for i in self.records:
			record = i.dump()
			interval = record.get('interval', self.default_interval)
			self.logger.debug("interval")
			self.logger.debug(interval)
			self.logger.debug("time")
			self.logger.debug(int(time.time()))
			self.logger.debug("timestamp")
			self.logger.debug(self.timestamp[record.get('_id')])
			self.logger.debug("diff")
			self.logger.debug(( int(time.time()) - self.timestamp[record.get('_id')]) )
			if ( int(interval) < ( int(time.time()) - self.timestamp[record.get('_id')]) ):
				self.logger.debug('ok');
				tfilter = json.loads(record.get('mfilter'))
				metric_list = self.manager.store.find(mfilter=tfilter)
				values = []
				i=1
				fn = record.get('type', False)
				if ( isinstance(fn, str)) :
					fn = [ fn ] 
				for metric in metric_list :
					m = metric.get('d')
					values.append( m[len(m)-1][1] ) 
					i = i + 1
				if ( fn ) :
					for i in fn :
						self.logger.debug(i)
						self.logger.debug(values)
						resultat = pyperfstore2.utils.aggregate_series(values, fn )
						self.logger.debug(resultat)	 
				self.timestamp[record.get('_id')] = int(time.time())
		
		
	def load (self, rec ) :
		record = rec.dump()
		rec.loaded = True
		self.storage.update(record.get('_id'), {'loaded': 'true' })
		if ( record.get('mfilter', False) ) :
			self.timestamp[record.get('_id')] = int(time.time())
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
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)

	def load_consolidation(self) :
		self.records = self.storage.find({ 'crecord_type': 'consolidation' }, namespace="object")
		for i in self.records :
			self.load(i)
				
			
			
			
	def unload_consolidation(self):
		record_list = self.storage.find({ 'crecord_type': 'consolidation' }, namespace="object")
		for i in record_list :
			self.storage.update(i._id, {'loaded': 'false' })
