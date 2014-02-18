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
from crecord import crecord
from cstorage import get_storage

import cevent
import logging
import threading

NAME='event_linker'

class engine(cengine):
	"""
		Event Linker engine.
		Use the event's RK to reference the event in a tree.
	"""

	def __init__(self, *args, **kwargs):
		super(engine, self).__init__(name=NAME, *args, **kwargs)

		self.account = caccount(user='root', group='root')
		self.storage = get_storage(namespace='events_trees', logging_level=logging.DEBUG, account=self.account)

		# We don't need to lock this structure, the only write operation done
		# is an append which is a thread-safe operation.
		self.trees = []

	def beat(self):
		"""
			Executed periodically, store the super-tree to the
			database.
		"""

		for tree in self.trees:
			record = self.storage.find_one(mfilter={'rk': tree['rk']})

			if not record:
				record = crecord(storage=self.storage, data=tree)

			else:
				record.data['rk'] = tree['rk']
				record.data['child_nodes'] = tree['child_nodes']

			record.save(storage=self.storage)

	def work(self, event, *args, **kwargs):
		"""
			Engine worker. Parse the event to add it to the super-tree.

			:param event: AMQP event.
			:type event: dict

			:returns: AMQP event as dict.
		"""

		rk = event['rk']
		rk_components = rk.split('.')

		current_rk = rk_components[0]
		current_node = {}

		# If the root node already exists, set the current node to it
		for tree in self.trees:
			if rk_components[0] == tree['rk']:
				current_node = tree
				break

		# Else, create a new node and set the current root node to it
		else:
			current_node['rk'] = rk_components[0]
			current_node['child_nodes'] = []

			self.trees.append(current_node)

		# Insert event's routing key in the tree
		for rkcomp in rk_components[1:]:
			current_rk = '{0}.{1}'.format(current_rk, rkcomp)

			# If the node already exists, just change the current node to the
			# one found
			for child in current_node['child_nodes']:
				if child['rk'] == current_rk:
					current_node = child
					break

			# If the node doesn't exist, then create it and change the current
			# node to the new one
			else:
				new_node = {
					'rk': current_rk,
					'child_nodes': []
				}

				current_node['child_nodes'].append(new_node)
				current_node = new_node

		return event
