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
from cstorage import get_storage
from caccount import caccount
from cselector import cselector

import logging
		
NAME="selector"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.selectors = []
		self.nb_beat = 0
		self.thd_warn_sec_per_evt = 1.5
		self.thd_crit_sec_per_evt = 2
		
	def pre_run(self):
		#load selectors
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))

		
	def beat(self):
		self.logger.debug('entered in selector BEAT')
		# Refresh selectors for work method
		
	
	def consume_dispatcher(self,  event, *args, **kargs):
		self.logger.debug('entered in selector consume dispatcher')
		# Gets crecord from amqp distribution
	
		selector = self.get_ready_record(event)
		if selector:

			event_id = event['_id']

			# Loads associated class
			selector = cselector(storage=self.storage, record=selector, logging_level=self.logging_level)

			self.logger.debug('%s found, start processing..' % (event_id))			
			# do I publish a selector event ? Yes if selector have to and it is time or we got to update status 
			if selector.dostate:
				try:
					#TODO improve this full mongo db request
					rk, selector_event = selector.event()
					self.logger.info('%s properly computed' % (event_id))		

				except Exception as e:
					self.logger.error('Unable to select all event matching this selector in order to publish worst state one form them. Exception : ' + str(e))
					event = None
				
				# Publish Sla information when available
				publishSla = selector.data.get('sla_rk', None)
				if publishSla:
					selector_event['sla_rk'] = publishSla
										
				# Ok then i have to update selector statement
				self.storage.update(event_id, {'state': selector_event['state']})
				self.amqp.publish(selector_event, rk, self.amqp.exchange_name_events)
				self.logger.debug("%s published event" % (selector.name))
					
			else:
				self.logger.debug('Nothing to do with this selector')
			
			#Update crecords informations	
			self.crecord_task_complete(event_id)
			
		self.nb_beat +=1
		#set record free for dispatcher engine


		
		

