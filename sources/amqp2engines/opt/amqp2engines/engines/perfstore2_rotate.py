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
import time

NAME="perfstore2_rotate"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)

		self.beat_interval=10
		
		self.kplan = "perfstore2:rotate:plan"

		self.rotation_interval = 60 * 60 * 24 # 24 hours
		self.key_by_beat = 200

		self.build_interval = 60 * 60 * 2 # 2 hours
		self.last_build = time.time()
		
	def pre_run(self):
		self.manager = pyperfstore2.manager()
		self.redis = self.manager.store.redis

		if not self.redis.exists(self.kplan):
			self.build_rotate_plan()

		self.beat()

	def build_rotate_plan(self):
		self.logger.info("Build rotate plan")
		start = time.time()

		self.last_build = start

		planned_keys = self.redis.zrange(self.kplan, 0, -1)
		keys = self.redis.keys()

		self.logger.info(" + Planned keys: %s", len(planned_keys))
		self.logger.info(" + Keys: %s", (len(keys) - 1)) # - self.kplan

		tmp_dict = {}
		for k in planned_keys:
			tmp_dict[k] = False

		to_add = [ k for k in keys if tmp_dict.get(k, True) and k != self.kplan ]

		tmp_dict = {}
		for k in keys:
			tmp_dict[k] = False

		to_rem = [ k for k in planned_keys if tmp_dict.get(k, True) and k != self.kplan ]

		rp = self.redis.pipeline()

		self.logger.info(" + Add %s keys", len(to_add))
		for key in to_add:
			rp.zadd(self.kplan, 0, key)

		rp.execute()

		self.logger.info(" + Remove %s keys", len(to_rem))
		for key in to_rem:
			rp.zrem(self.kplan, key)

		rp.execute()

		elapsed = (time.time() - start) * 1000
		self.logger.info("Done in %.2f ms", elapsed)

	def beat(self):
		self.logger.debug("Start rotation")
		start = time.time()

		if (self.last_build + self.build_interval < start):
			self.build_rotate_plan()
			return

		rp = self.redis.pipeline()

		keys = self.redis.zrangebyscore(self.kplan, 0, int(start), start=0, num=self.key_by_beat)

		## Set net time
		for key in keys:
			rp.zadd(self.kplan, int(start + self.rotation_interval), key)
		rp.execute()

		self.logger.debug(" + Keys: %s" % len(keys))

		## Work
		for key in keys:
			self.manager.rotate(_id=key)
	
		elapsed = (time.time() - start)
		self.counter_event += len(keys)
		self.counter_worktime += elapsed

		if elapsed > self.beat_interval - 3:
			self.logger.warning("Rotation time %s s is to close from beat interval (%s s)" % (int(elapsed), self.beat_interval) )

		self.logger.debug("Done in %.2f ms", int(elapsed*1000))
