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

import unittest
import time

import logging

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(name)s %(levelname)s %(message)s',
                    )

from caccount import caccount
from cstorage import get_storage
from cdowntime import Cdowntime

root_account = caccount(user="root", group="root")
storage = get_storage(account=root_account , namespace='unittest', logging_level=logging.DEBUG)


class KnownValues(unittest.TestCase): 
	def setUp(self):
		self.backend = storage.get_backend('unittest')
		self.cdowntime = Cdowntime(storage)
		#Overidding default backend
		self.cdowntime.backend = self.backend

	def test_01_Method_get_filter_no_data(self):
		mongo_filter = self.cdowntime.get_filter()
		if mongo_filter:
			raise()

	def test_02_Method_get_filter_data_feed(self):

		self.backend.insert({
			'component'	: 'component_test_1',
			'resource'	: 'resource_test_1',
			'type'		: 'downtime',
			'start'		: 0,
			'end'		: time.time() + 10000
		})
		self.backend.insert({
			'component'	: 'component_test_2',
			'resource'	: 'resource_test_2',
			'type'		: 'downtime',
			'start'		: time.time() + 10000,
			'end'		: time.time() + 10000
		})

		mongo_filter = self.cdowntime.get_filter()
		if not mongo_filter:
			raise('Should have selected something')
		if len(mongo_filter['$and']) != 1:
			raise('filter should be defined for excactly one element')
		for element in mongo_filter['$and'][0]['$and']:
			for element_type in ['component', 'resource']:
				if element_type in element['$ne'] and element['$ne'][element_type] not in  ['component_test_1', 'resource_test_1']:
					raise('iterated keys should had either component_test_1|resource_test_1 values')

	def test_99_DropNamespace(self):
		storage.drop_namespace('unittest')


if __name__ == "__main__":
	unittest.main(verbosity=2)


