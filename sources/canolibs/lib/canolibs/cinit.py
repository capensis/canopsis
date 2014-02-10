#!/usr/bin/env python
# --------------------------------
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

import signal, time

class cinit(object):
	class getHandler(object):
		def __init__(self, logger):
			self.logger = logger
			self.RUN = True

		def status(self):
			return self.RUN

		def signal_handler(self, signum, frame):
			self.logger.warning("Receive signal to stop daemon...")
			if self.callback:
				self.callback()
			self.stop()

		def run(self, callback=None):
			self.callback = callback
			signal.signal(signal.SIGINT, self.signal_handler)
			signal.signal(signal.SIGTERM, self.signal_handler)
			
		def stop(self):
			self.RUN = False

		def set(self, statut):
			self.RUN = statut

		def wait(self):
			while self.RUN:
				try:
					time.sleep(1)
				except:
					break
			self.stop()
