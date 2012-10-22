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

import time
import random
import logging

import pyperfstore

from pyperfstore import node
from pyperfstore import metric
from pyperfstore import dca

from pyperfstore import filestore
from pyperfstore import memstore
from pyperfstore import mongostore

logging.basicConfig(level=logging.INFO,
	format='%(name)s %(levelname)s %(message)s',
)

rotate_plan = {
	'PLAIN': 0,
	'TSC': 3,
}

interval = 300
day = 30
point_per_dca=None #auto
#point_per_dca=300

def bench_store(store, interval=60, duration=60*60*24, point_per_dca=None):
	print "Start Bench ..."
	
	msize = store.size()
	
	mynode = node('nagios.Central.check.service.localhost', storage=store, rotate_plan=rotate_plan, point_per_dca=point_per_dca)

	# 1 value / 5 min = 8928 values/month = 107136 values/year
	timestamp = int(time.time())
	#timestamp = 1
	bench_start = timestamp
	
	nb = duration / interval
	
	start = time.time()
	for i in range(1,nb+1):
		mynode.metric_push_value(dn='state', unit=None, value=1, timestamp=timestamp)
		#mynode.metric_push_value(dn='state-warning', unit=None, value=0, timestamp=timestamp)
		#mynode.metric_push_value(dn='state-critical', unit=None, value=0, timestamp=timestamp)
		mynode.metric_push_value(dn='state-downtime', unit=None, value=0, timestamp=timestamp)
		value = random.random() * 10
		mynode.metric_push_value(dn='load1', unit=None, value=value, timestamp=timestamp)
		mynode.metric_push_value(dn='load5', unit=None, value=value, timestamp=timestamp)
		mynode.metric_push_value(dn='load15', unit=None, value=value, timestamp=timestamp)

		timestamp += interval

	bench_stop = timestamp

	nb = nb * 3
	elapsed = time.time() - start
	
	print " + WRITE:"
	print "    + %.2f days" % (duration / (60*60*24))
	msize = store.size()
	print "    + %.2f MB (%.2f MB/Year)" % ((msize / 1024.0 / 1024.0), ((msize / float(duration))/ 1024.0 / 1024.0) *  60*60*24*365)
	#nsize = mynode.size()
	#print "    + %.2f MB (%.2f MB/Year)" % ((nsize / 1024.0 / 1024.0), ((nsize / float(duration))/ 1024.0 / 1024.0) *  60*60*24*365)
	#print "    + Delta: %s B" % (msize - nsize)
	print "    + %s values in %s seconds" % ( nb, elapsed)
	print "    + %s values per second" % (int(nb/elapsed))
	print ""

	start = time.time()
	mynode = node('nagios.Central.check.service.localhost', storage=store)
	print "Get values between %s and %s" % (bench_start, bench_stop)
	values = mynode.metric_get_values(dn='load1', tstart=bench_start, tstop=bench_stop)
	nb = len(values)
	elapsed = time.time() - start
	print " + READ:"
	print "    + %s values in %s seconds" % ( nb, elapsed)
	print "    + %s values per second" % (int(nb/elapsed))
	print ""
	
	mynode.pretty_print()
	size = store.size()
	
	start = time.time()
	mynode.remove()
	elapsed = time.time() - start
	
	size = store.size()
	print " + REMOVE:"
	print "    + %.2f MB" % (size / 1024.0 / 1024.0)
	print "    + %s seconds" % elapsed
	print ""


print "Mongo Store"
storage = mongostore(mongo_safe=False, mongo_collection='pyperfstorebench')
storage.drop()

bench_store(	storage,
				interval=interval,
				duration=60*60*24*day,
				point_per_dca=point_per_dca)

#print "Files Store"
#storage = filestore(base_path="/tmp/data")
#bench_store(storage)

#print "Memory Store"
#storage = memstore()
#bench_store(storage)
