#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

import unittest, sys,json
import logging
import time
import random

logging.basicConfig(level=logging.INFO,
	format='%(name)s %(levelname)s %(message)s',
)

import pyperfstore2
manager = None
name = 'nagios.Central.check.service.localhost9.ping'
component = 'localhost9'
resource = 'ping'

#start = int(time.time())
start = 0
ut_start = start

nb = 0
n=60
stop = n

meta_data = {'dn': (component, resource), 'retention': 2 }

class KnownValues(unittest.TestCase):
	def setUp(self):
		pass

	def test_01_Init(self):
		global manager

		manager = pyperfstore2.manager(
			mongo_collection='unittest_perfdata2',
			dca_min_length=50,
			logging_level=logging.DEBUG,
			redis_db=1)
		
		manager.store.drop()

	def test_02_Push(self, set_timestamp=True):
		global start, stop, nb
		
		if set_timestamp:
			manager.timestamp = start
		
		for i in range(start+1, stop+1):
			manager.push(name=name, value=i, timestamp=i, meta_data=meta_data)
			nb+=1
		
		manager.store.sync()

		start = i
		stop = start + n
		manager.timestamp = start
					
	def test_03_Rotate(self):
		manager.rotateAll()

	def test_04_Push(self):
		self.test_02_Push()
		
	def test_05_Functions(self):
		data = manager.find(name=name)
		if data.count() != 1:
			raise Exception('Invalid meta count')
			
		if len(data[0].get('c', [])) != 1:
			raise Exception('Invalid rotation')
			
		data = manager.find(name=name, limit=1, data=False)
		if data.get('d', None):
			raise Exception('Data field is present')

			
	def test_07_Get_points(self):
		points = manager.get_points(name=name, tstart=ut_start, tstop=stop)
		
		print "Total: %s" % nb
		print "Nb points: %s" % len(points)
		if len(points) != nb:
			raise Exception('Invalid points count')
				
		points = pyperfstore2.utils.aggregate(points, max_points=50, atype='MEAN', mode='by_point')
		print "Nb aggregate points: %s" % len(points)
		
		points = manager.get_points(name=name, tstart=100, tstop=119)
		if len(points) != 20:
			raise Exception('Invalid points count: %s' % len(points))

	def test_08_prev_next_points(self):
		my_start = ut_start+75
		my_stop = stop-75
		points = manager.get_points(name=name, tstart=my_start, tstop=my_stop)
		print "Nb points: %s" % len(points)
		my_nb = len(points)
		
		points = manager.get_points(name=name, tstart=my_start, tstop=my_stop, add_next_point=True)
		if len(points) != my_nb + 1:
			raise Exception('Invalid add_next_point (%s)' % len(points))
			
		points = manager.get_points(name=name, tstart=my_start, tstop=my_stop, add_prev_point=True)
		if len(points) != my_nb + 1:
			raise Exception('Invalid add_prev_point (%s)' % len(points))
			
		points = manager.get_points(name=name, tstart=my_start, tstop=my_stop, add_next_point=True, add_prev_point=True)
		if len(points) != my_nb + 2:
			raise Exception('Invalid add_prev_point +  add_prev_point (%s)' % len(points))
			
			
	def test_09_Get_point(self):
		point = manager.get_last_point(name=name)
		print "Point: %s" % point
		
		if len(point) != 2:
			raise Exception('Invalid get_last_point (%s)' % point)
			
		if point[0] != (stop-n):
			raise Exception('Invalid last point timestamp (%s != %s)' % (point[0], stop-n))
			
		point = manager.get_point(name=name, ts=ut_start)
		print "Point: %s" % point
		
		if len(point) != 2 and point[0] != ut_start:
			raise Exception('Invalid get_point (%s -> %s)' % (ut_start, point))
	
	def test_10_ShowAll(self):
		manager.showStats()
		manager.showAll()
				
	def test_11_clean(self):
		## Rotate	
		manager.rotate(name=name)
		
		## Push n points
		self.test_02_Push()
		
		## TODO: re-code this function !
		## Clean old dca
		#cleaned = manager.clean(name=name, timestamp=ut_start+n)
		#if not cleaned:
		#	raise Exception('Data must be cleaning: %s' % cleaned)
			
		cleaned = manager.cleanAll(timestamp=ut_start+n)
		if cleaned:
			raise Exception('Data cleaning: %s' % cleaned)	
			
		## Show stats
		manager.showAll()			
		
		points = manager.get_points(name=name, tstart=ut_start, tstop=stop)
		if len(points) != 120:
			raise Exception('Invalid count %s' % len(points))	
		
	def test_97_Remove(self):
		manager.remove(name=name)
		meta = manager.get_meta(name=name)
		
		if meta:
			raise Exception('Impossible to delete')
			
	def test_98_Add_Data_after_purge(self):
		_id = manager.get_id(name=name)
								
		## Push n points
		self.test_02_Push()
		
		dca = manager.get_meta(_id)
		
		if not dca:
			raise Exception('DCA not upserted')
		
		points = manager.get_points(name=name, tstart=ut_start, tstop=stop)
		
		if not points:
			raise Exception('Point not inserted')
		else:	
			if len(points) != 60:
				raise Exception('Invalid count: %s ' % len(points))
	
	def test_99_Drop(self):
		manager.showAll()
		manager.store.drop()
	
if __name__ == "__main__":
	unittest.main(verbosity=2)
