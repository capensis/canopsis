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
from cstorage import get_storage
from caccount import caccount

import pyperfstore2
import logging
import time

NAME="perfstore2_rotate"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME,logging_level=logging.DEBUG, *args, **kargs)

		self.beat_interval=10
		
		self.kplan = "perfstore2:rotate:plan"

		self.rotation_interval = 60 *10#* 60 * 24 # 24 hours


		self.last_build = time.time()
		
	def pre_run(self):
		self.manager = pyperfstore2.manager(logging_level=logging.DEBUG)

		self.beat()

	def beat(self):
		self.logger.debug("Start rotation")
		start = time.time()

		metric_to_rotate = self.manager.store.daily_collection.find({'insert_date': {'$lte': start - self.rotation_interval}})

		metric_count = 0
		for metric in metric_to_rotate:
			self.manager.rotate(metric['_id'], metric['values'])
			metric_count += 1


		
		elapsed = time.time() - start
		self.counter_event += metric_count
		self.counter_worktime += elapsed

		if elapsed > self.beat_interval - 3:
			self.logger.warning("Rotation time %s s is to close from beat interval (%s s)" % (int(elapsed), self.beat_interval) )

		self.logger.debug("Done in %.2f ms", int(elapsed*1000))
