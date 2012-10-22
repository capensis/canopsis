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

import unittest, logging
import time
from ctimer import ctimer

mytimer = None

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_1_Init(self):
		global mytimer
		mytimer = ctimer(logging_level=logging.DEBUG)

	def test_2_Start(self):
		mytimer.start()

	def test_3_Stop(self):
		time.sleep(1)
		mytimer.stop()
		if not (mytimer.elapsed > 0.9 and mytimer.elapsed < 1.1):
			raise Exception('Invalid elapsed time ...')

	def test_4_Task(self):
		def task(_id="defaultid", name='defaultname'):
			print time.time(), _id, name
			time.sleep(0.7)

		start = time.time()
		mytimer.start_task(task=task, interval=1, count=3, _id='myid', name='myname')
		stop = time.time()
		elaps = round(stop - start, 1)
		print "Start:", start, "Stop:", stop, "Elapsed:", elaps

		if elaps != 2.7:
			raise Exception('Invalid elapsed time ... (%s != 2.7)' % elaps)

		
if __name__ == "__main__":
	unittest.main(verbosity=2)
	
