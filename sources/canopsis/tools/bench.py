#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
from pyperfstore3.manager import Manager

logging.basicConfig(
	level=logging.INFO, format='%(name)s %(levelname)s %(message)s')

interval = 300
day = 30
duration = 60 * 60 * 24 * day

name = 'bench'
DATA_ID = name
manager = Manager(data_name=name)
manager.remove(data_id=DATA_ID, with_meta=True)


def bench_store(interval=interval, duration=60 * 60 * 24):
	print("Start Bench ...")

	msize = manager.store.size()

	# 1 value / 5 min = 8928 values/month = 107136 values/year
	timestamp = int(time.time())
	bench_start = timestamp

	nb = duration / interval
	print(" + write %s loop" % nb)

	start = time.time()
	for i in range(1, nb + 1):
		manager.push(name="%s%s" % (name, 'state'), value=1, timestamp=timestamp)
		manager.push(
			name="%s%s" % (name, 'state-downtime'), value=1, timestamp=timestamp)

		value = random.random() * 10
		manager.push(name="%s%s" % (name, 'load1'), value=value, timestamp=timestamp)
		manager.push(name="%s%s" % (name, 'load5'), value=value, timestamp=timestamp)
		manager.push(
			name="%s%s" % (name, 'load15'), value=value, timestamp=timestamp)

		timestamp += interval

		mod = (i * interval) % 86400
		if mod == 0:
			manager.midnight = timestamp
			manager.rotateAll()

	bench_stop = timestamp

	nb = nb * 3
	elapsed = time.time() - start

	print(" + WRITE:")
	print("    + %.2f days" % (duration / (60 * 60 * 24)))
	msize = manager.store.size()
	print("    + %.2f MB (%.2f MB/Year)" % (
		(msize / 1024.0 / 1024.0),
		((msize / float(duration)) / 1024.0 / 1024.0) * 60 * 60 * 24 * 365))
	#nsize = mynode.size()
	#print "    + %.2f MB (%.2f MB/Year)" % ((nsize / 1024.0 / 1024.0),
	# ((nsize / float(duration))/ 1024.0 / 1024.0) *  60*60*24*365)
	#print "    + Delta: %s B" % (msize - nsize)
	print("    + %s values in %s seconds" % (nb, elapsed))
	print("    + %s values per second" % (int(nb / elapsed)))
	print("")

	nb = 0
	start = time.time()
	print("Get values between %s and %s" % (bench_start, bench_stop))
	values = manager.get_points(
		name="%s%s" % (name, 'load1'), tstart=bench_start, tstop=bench_stop)
	nb += len(values)
	values = manager.get_points(
		name="%s%s" % (name, 'load5'), tstart=bench_start, tstop=bench_stop)
	nb += len(values)
	values = manager.get_points(
		name="%s%s" % (name, 'load15'), tstart=bench_start, tstop=bench_stop)
	nb += len(values)
	elapsed = time.time() - start
	print(" + READ:")
	print("    + %s values in %s seconds" % (nb, elapsed))
	print("    + %s values per second" % (int(nb / elapsed)))
	print("")


bench_store(interval=interval, duration=60 * 60 * 24 * day)

manager.remove(data_id=DATA_ID)
