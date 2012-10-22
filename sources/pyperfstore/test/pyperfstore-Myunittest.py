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

import unittest, sys,json
import logging
import time
import random

import pyperfstore

from pyperfstore import filestore
from pyperfstore import memstore
from pyperfstore import node
from pyperfstore import metric
from pyperfstore import dca
from pyperfstore import pmath

logging.basicConfig(level=logging.INFO,
	format='%(name)s %(levelname)s %(message)s',
)

mynode = None
storage = None
timestamp = None
refvalues = [ ]
refvalues.append([0, 0])

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_Init(self):
		global mynode, storage, timestamp
		#storage = filestore(base_path="/tmp/")
		storage = memstore()
		_id = 'nagios.Central.check.service.localhost9'
		mynode = node(_id=_id, dn=_id, point_per_dca=100, storage=storage)

		timestamp = 1

	def test_02_PushValue(self):
		global timestamp, refvalues
		# 1 value / 5 min = 8928 values/month = 107136 values/year
		interval = 1
		nb = 1000
		for i in range(1,nb):
			
			value = random.random()
			mynode.metric_push_value(dn='load1', value=value, timestamp=timestamp)
			refvalues.append([timestamp, value])

			mynode.metric_push_value(dn='load5', value=random.random(), timestamp=timestamp)
			mynode.metric_push_value(dn='load15', value=random.random(), timestamp=timestamp)

			timestamp += interval

		mynode.pretty_print()

	def test_03_Load(self):
		global mynode

		dump1 = mynode.dump()

		metric1 =  mynode.metric_dump(dn='load1')

		del mynode

		mynode = node('nagios.Central.check.service.localhost9', storage=storage)	

		dump2 = mynode.dump()
		metric2 =  mynode.metric_dump(dn='load1')

		del dump1['writetime']
		del dump2['writetime']

		del metric1['writetime']
		del metric2['writetime']

		if dump1 != dump2:
			print "First:"
			print dump1
			print "Second:"
			print dump2
			raise Exception('Invalid load of nodes...')

		if metric1 != metric2:
			print "First:"
			print metric1
			print "Second:"
			print metric2
			raise Exception('Invalid load of metrics ...')

	def test_04_Second_PushValue(self):
		self.test_02_PushValue()

		print ""
		mynode.pretty_print()
		pass


	def test_05_GetBy(self):
		_id =  mynode.metric_get_id(dn='load1')

	def test_06_GetValues(self):
		last = timestamp - 1

		print "Last: %s" % last

		## Get first 100 values
		start = time.time()
		values = mynode.metric_get_values(dn='load1', tstart=1, tstop=100)
		print " + %s Old values in %s ms" % (len(values), ((time.time() - start) * 1000))

		if len(values) != 100:
			print "Count: %s" % len(values)
			raise Exception('Invalid Old count (len: %s)' % len(values))

		if values != refvalues[1:100+1]:
			print values
			print refvalues[1:100+1]
			raise Exception('Invalid Old Data')


		## Get last 100 values
		start = time.time()
		values = mynode.metric_get_values(dn='load1', tstart=last-99, tstop=last)
		print " + %s Recent values in %s ms" % (len(values), ((time.time() - start) * 1000))

		if len(values) != 100:
			print "Count: %s" % len(values)
			raise Exception('Invalid Recent count (len: %s)' % len(values))

		
		if values != refvalues[last-99:last+1]:
			print values
			print refvalues[last-99:last+1]
			raise Exception('Invalid Recent Data')

		## Get middle 100 values
		start = time.time()
		values = mynode.metric_get_values(dn='load1', tstart=last-499, tstop=last-400)
		print " + %s Middle values in %s ms" % (len(values), ((time.time() - start) * 1000))

		if len(values) != 100:
			print "Count: %s" % len(values)
			raise Exception('Invalid Middle count (len: %s)' % len(values))

		
		if values != refvalues[last-499:last-400+1]:
			print values
			print refvalues[last-499:last-400+1]
			raise Exception('Invalid Middle Data')


	def test_07_aggregate(self):
		ori_values = mynode.metric_get_values(dn='load1', tstart=1, tstop=100)
		values = pyperfstore.pmath.aggregate(ori_values, max_points=50)

		if len(values) != 50:
			raise Exception('Invalid aggregate by points (len: %s)' % len(values))
			
		values = pyperfstore.pmath.aggregate(ori_values, time_interval=10, mode='by_interval')
		if len(values) != 10:
			raise Exception('Invalid aggregate by interval (len: %s)' % len(values))		
		

	#def test_08_candlestick(self):
	#	##### DRAFT !
	#	values = mynode.metric_get_values(dn='load1', tstart=1, tstop=1000, aggregate=False)

	#	values = pyperfstore.pmath.candlestick(values, window=100)

		#if len(values) != 10:
		#	raise Exception('Invalid candlestick (len: %s)' % len(values))

	def test_09_timesplit(self):
		##### DRAFT !
		values = mynode.metric_get_values(dn='load1', tstart=1, tstop=1000, aggregate=False)

		start = time.time()
		pyperfstore.pmath.timesplit(values, 35, 632)
		print " + Split %s points %s ms" % ( len(values), ((time.time() - start) * 1000))


	def test_10_fill_interval(self):
		values = [
			[0, 0],
			[10, 1],
			[20, 2],
			[30, 3],
			[90, 9],
		]
		values = pmath.fill_interval(values, 10)
		
		if values != [[0, 0], [10, 1], [20, 2], [30, 3], [40, 3], [50, 3], [60, 3], [70, 3], [80, 3], [90, 9]]:
			raise Exception('Invalid fill, %s' % values)
			
	def test_11_get_all_nodes_and_metrics(self):
		records = storage.get_all_nodes()
		if not records:
			raise Exception('Impossible to get nodes')
		print('')
		print('nodes are :')
		for r in records:
			print(r)
		
		records = storage.get_all_metrics()
		if not records:
			raise Exception('Impossible to get metrics')
		print('')
		print('metrics are :')
		for r in records:
			print(r)
		print('')
			
	def test_90_Functions(self):
		dn = mynode.metric_get_all_dn()
		print "dn:", dn
		if not dn:
			raise Exception('Impossible to get dn')

	def test_99_Remove(self):
		global mynode

		mynode.metric_remove('load1')
		
		values = mynode.metric_get_values(dn='load1', tstart=1, tstop=100)
		if values:
			raise Exception('Impossible to remove "load1"')
			
		values = mynode.metric_get_values(dn='load5', tstart=1, tstop=100)
		if not values:
			raise Exception('"load5" removed ?')
		
		mynode.metric_remove_all()
		
		values = mynode.metric_get_values(dn='load5', tstart=1, tstop=100)
		if values:
			raise Exception('Impossible to remove 	all')

		del mynode
		
		



if __name__ == "__main__":
	unittest.main(verbosity=2)
