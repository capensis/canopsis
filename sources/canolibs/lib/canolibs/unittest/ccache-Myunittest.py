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

from ccache import ccache

from caccount import caccount
from cstorage import cstorage

import time

STORAGE = None
CACHE = None

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_Init(self):
		global CACHE
		CACHE = ccache(STORAGE)

	def test_02_Put(self):
		CACHE.put('mycache1', {'data': 'titit'})

	def test_03_Get(self):
		time.sleep(2)
		data = CACHE.get('mycache1', 1)
		if data:
			raise Exception('Data must be expired ...')

		CACHE.put('mycache2', {'data': 'titit'})
		time.sleep(2)
		data = CACHE.get('mycache2', 5)

		if not data:
			raise Exception('Data must be Ok ...')

	def test_04_Get(self):
		CACHE.put('mycache3', {'data': 'titit'})
		CACHE.remove('mycache3')

	def test_05_Decorator(self):
	
		@CACHE.deco('toto', 5, ['i'])
		def test(i):
			return int(time.time())

		t1 = test(i=1)
		time.sleep(2)
		t2 = test(i=1)
		
		if t1 != t2:
			raise Exception('Error in decorator ...')

	def test_99_DropNamespace(self):
		STORAGE.drop_namespace('unittest')
		pass

if __name__ == "__main__":
	STORAGE = cstorage(caccount(user="root", group="root"), namespace='unittest')
	unittest.main(verbosity=1)
	


