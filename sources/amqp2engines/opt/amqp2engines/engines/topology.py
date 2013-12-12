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
import cevent
import logging

from caccount import caccount
from cstorage import get_storage

import time
from datetime import datetime
from random import random

NAME="topology"

import sys, os
sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/%s/' % NAME))

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
						
		self.beat_interval = 60
		self.nb_beat = 0
		# Operator cache
		self.modules = {}
		
		# All ids in all topos

		self.stateById = {}
		self.topos = []
		
		# Beat
		self.doBeat = False
		self.normal_beat_interval = 300
		self.lastBeat = int(time.time()) - self.normal_beat_interval		
		
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		
		# TODO: Not compatible with HA !
		self.topo_unloadAll()
		
		self.beat()

	def topo_load(self):
		count = self.storage.count({'crecord_type': 'topology', 'enable': True, 'loaded': False}, namespace="object")
		if count:
			self.logger.debug("Load topologies")
			self.topos = self.storage.find({'crecord_type': 'topology', 'enable': True, 'loaded': False}, namespace="object", raw=True)
				
	def topo_unload(self):
		self.logger.debug("Unload topologies")
		for topo in self.topos:
			self.storage.update(topo['_id'], {'loaded': False})
		
		
	def topo_unloadAll(self):
		records = self.storage.find({'crecord_type': 'topology', 'enable': True},  mfields=['_id'], namespace="object", raw=True)
		for record in records:
			self.storage.update(record['_id'], {'loaded': False})
			

		del self.topos
		
		self.topos = []
	
	def default_Operator_fn(self, states, options={}):
		self.logger.debug("Calcul state of sample: %s" % states)
		
		if isinstance(states, int):
			return states
		elif isinstance(states, list):
			if len(states) == 0:
				self.logger.warning("No state in 'states'")
				return 0
			elif len(states) == 1:
				return states[0]
			else:
				return 0
		else:
			self.logger.warning("Invalid states type: %s" % type(states))
			return 0
	
	def topo_getOperator_fn(self, name):
		module = self.modules.get(name, None)
		if not module:
			try:
				module = __import__(name)
				module.logger = self.logger
				self.modules[name] = module
				return module.operator
			except Exception, err:
				self.logger.warning("Impossible to load operator '%s' (%s)" % (name, err))
				return self.calcul_state
		else:
			return module.operator	

		
	def topo_fillChilds(self, topo):
		for conn in topo['conns']:
			srcId = conn[0]
			dstId = conn[1]
			topo['nodes'][dstId]['childs'].append(srcId)
					
	def topo_getState(self, topo):
		state = 0
		state_type = 1
		rootNode = topo['nodes'][topo['root']]
		
		### Recursive function
		def parseChilds(parent, level=0):
			
			_id  = parent['_id']
			state = self.stateById.get(_id, {})
			childIds = parent.get('childs', [])
			
			## For DEBUG
			prefix = ""
			for i in range(0, level):
				prefix = "  " + prefix
			
			self.logger.debug("%s|-> %s (State: %s, type: %s)" % (prefix, _id , state.get('state', None), state.get('state_type', None)))
			
			states = []
			
			# Browse all tree
			if childIds:
				# Parse childs
				for childId in childIds:
					state = parseChilds(topo['nodes'][childId], level+1)
					states.append(state)
				
				# Calcule state of parent
				try:
					state = parent['calcul_state'](states=states, options=parent.get("options", {}))
				except Exception, err:
					self.logger.error("Impossible to calcul state of %s (%s)" % (_id, err))
					state = self.default_Operator_fn(states=states)
				
				# Set state
				self.stateById[_id] = {'state': state, 'state_type': 1, 'previous_state': state}
				
				return state 
			else:
				## End of tree, return parent state
				if not state:
					self.logger.error("No state for %s" % _id )
					return 3
				
				if not state['state_type']:
					return state['previous_state']
				else:
					return state['state']
		
		## State recursivity
		state = parseChilds(rootNode)

		return { 'state': state, 'state_type': state_type }

	def topo_dump4Ui(self, topo):

		def parseChilds(parent):
			label = parent.get('label', None)
			if not label:
				if  parent.get('resource', None):
					label = "%s %s" % (parent.get('component', ''), parent.get('resource', ''))
				else:
					label = parent.get('component', '')
			
			childs = []
			if parent.get('childs', []):
				for child in parent['childs']:
					childs.append(parseChilds(topo['nodes'][child]))
			
			states = self.stateById.get(parent['_id'], { 'state': 3, 'state_type': 1, 'previous_state': 3 })
			
			return {
				'_id': "%s-%s" % (parent['_id'], int(random() * 10000)),
				'event_id': parent['_id'],
				'name': label,
				'childs': childs,
				'state': states['state'],
				'state_type': states['state_type'],
				'previous_state': states['previous_state']
			}	
		

		root = topo['nodes'][topo['root']]
		tree = parseChilds(root)
		return tree
	
	def post_run(self):
		self.topo_unload()
		
			
	def beat(self):
		self.logger.debug('entered in topology BEAT')
		self.topo_load()
		# Refresh selectors for work method
		self.nb_beat += 1	
		
	
	def consume_dispatcher(self,  event, *args, **kargs):
	
		self.logger.debug('entered in topology consume dispatcher')
		# Gets crecord from amqp distribution
		topology = self.get_ready_record(event)
		
		if topology:	
			event_id = event['_id']
			ids = []
			
			# Parse topo
			for topo in self.topos:				
				self.logger.debug("Parse topo '%s': %s Nodes with %s Conns" % (topo['crecord_name'], len(topo['nodes']), len(topo['conns'])))
									
				nodes = topo['nodes']
				topo['ids'] = []
				for key in nodes:
					_id = nodes[key].get('_id', None)
					if _id and _id  not in ids:
						topo['ids'].append(_id)
						if _id  not in ids:
							ids.append(_id)
				
		
				topo['nodesById'] = {}
				
				for key in topo['nodes']:
					node = topo['nodes'][key]
					
					_id = node['_id']
					
					if not node.get('calcul_state', None):
						if node.get('event_type', None) == 'operator':
							node['calcul_state'] = self.topo_getOperator_fn(_id)
							_id = "%s-%s" % (_id, int(random() * 10000))
							node['_id'] = _id
						else:
							node['calcul_state'] = self.default_Operator_fn
						
					topo['nodesById'][_id] = node
					node['childs'] = []
							
				self.logger.debug("Fill node's childs")
				self.topo_fillChilds(topo)
			
		
			# Get all states of all topos
			self.stateById = {}
			records = self.storage.find(mfilter={'_id': {'$in': ids}}, mfields=['state', 'state_type', 'previous_state'], namespace='events')
			for record in records:
				self.stateById[record['_id']] = {
					'state': record['state'],
					'state_type': record.get('state_type', 1),
					'previous_state': record.get('previous_state', record['state'])
				}
		

			## Parse tree for calcul state
			self.logger.debug(" + Calcul state:")
			states_info = self.topo_getState(topo)

			self.logger.debug("'%s': State: %s" % (topo['crecord_name'], states_info))
			self.storage.update(topo['_id'], {'state': states_info['state']})
			
			event = cevent.forger(
				connector =			NAME,
				connector_name =	"engine",
				event_type =		"topology",
				source_type =		"component",
				component =			topo['crecord_name'],
				state =				states_info['state'],
				state_type =		states_info['state_type'],
				output =			"",
				long_output =		"",
				#perf_data =			None,
				#perf_data_array =	[],
				display_name =		topo.get('display_name', None)
			)
			
			# Extra fields			
			event['nestedTree'] = self.topo_dump4Ui(topo)
	

			rk = cevent.get_routingkey(event)

			self.logger.debug("Publish event on %s" % rk)
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)
			self.crecord_task_complete(event_id)

			
		else:
			self.logger.warning('topology not able to load crecord properly, topology not threaten.')		
	
	
	def work(self, event, *args, **kargs):
		if not self.doBeat:
			for topo in self.topos:
				if  event['rk'] in topo['ids']:
					self.doBeat = True
					break
