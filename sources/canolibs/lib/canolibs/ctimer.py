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
import clogging


class ctimer(object):
	def __init__(self):
		self.started = False
		self.logger = clogging.getLogger()
		self.RUN = True

	def start(self):
		#self.logger.debug("Start timer")
		self.started = True
		self.starttime = time.time()

	def stop(self):
		if self.started:
			#self.logger.debug("Stop timer")
			self.endtime = time.time()
			self.elapsed = self.endtime - self.starttime
			self.logger.debug("Elapsed time: %f ms", self.elapsed * 1000)

		self.started = False

	def start_task(self, task, interval=1, count=None, *args, **kargs):
		i=0
		tcount = 0
		start = time.time()
		self.logger.debug("Start task ...")
		while self.RUN:
			task(*args, **kargs)
			if count:
				tcount +=1
				if tcount >= count:
					break

			derive = time.time() - (start + (i*interval))
			i+=1
			if i >= 100:
				start = start + (i*interval)
				i=0

			pause = ((start + (i*interval)) - time.time())
			if pause < 0:
				pause = 0

			self.logger.debug("i: %s, Start: %s, Derive: %s, Pause: %s" % (i, start, derive, pause))
			try:
				step = int(pause / 0.5)
				rest = pause - (step * 0.5)
				self.logger.debug(" + Sleep: %s seconds (%s * 0.5)" % ((step * 0.5), step))
				if step > 0:
					for x in range(step-1):
						if self.RUN:
							time.sleep(0.5)
						else:
							self.logger.debug(" + Break !")
							break
				if self.RUN:
					self.logger.debug(" + Sleep %s seconds" % rest)
					pause = ((start + (i*interval)) - time.time())		
					time.sleep(pause)
			except:
				self.logger.debug(" + Exception !")
				self.RUN = False

		self.logger.debug("End of task ...")

	def stop_task(self):
		if self.RUN:
			self.logger.debug("Stop task ...")
			self.RUN = False

	def __del__(self):
		self.stop_task()
		
