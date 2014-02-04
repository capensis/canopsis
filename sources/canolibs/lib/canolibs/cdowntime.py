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

from crecord import crecord

import time



class Cdowntime(crecord):
	"""
		This class provide easy management for downtime by allowing component/resource test against any downtime at now time
	"""
	def __init__(self, storage):
		self.storage = storage
		self.backend = storage.get_backend('entities')


	def reload(self, delta_beat=0):
		""" Loads current downtimes being active
				delta_beat takes care of engines beat interval. for accurate measure,
				it should be equal to 0
		"""
		now = time.time()
		query = {'type': 'downtime', '$and':
			[
				{'start': { '$lte': now - delta_beat}},
				{'end'	: { '$gte': now +  delta_beat}}
			]
		}
		downtimes = self.backend.find(query)
		self.downtimes = [downtime for downtime in downtimes]
		return self.downtimes


	def is_downtime(self, component, resource):
		#Tests whether or not given component/resource information exists in downtime list.
		#If any, it s downtime and engines should operate consequently

		now = time.time()
		for downtime in self.downtimes:
			if downtime['component'] == component and downtime['resource'] == resource and downtime['start'] < now and downtime['end'] > now:
				return True
		return False

	def get_filter(self, mini=False):

		""" Builds a mongodb filter for downtime exclustion. mini is used for metric queries."""

		self.reload()

		if mini:
			component 	= 'co'
			resource 	= 're'
		else:
			component 	= 'component'
			resource 	= 'resource'


		if not self.downtimes:
			return None

		new_field = []
		for downtime in self.downtimes:
			new_filter = [
				{'$ne'	: {component: downtime['component']}}, 
				{'$ne'	: {resource: downtime['resource']}}
			]
			if mini:
				new_field += new_filter
			else:
				new_field.append({'$and': new_filter})


		if mini:
			return new_field
		else:
			return {'$and': new_field}
