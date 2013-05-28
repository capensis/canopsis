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
import time, json, logging
from bson import BSON

from camqp import camqp
from cinit import cinit
## Engines path
import sys, os
sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))

## Configurations

DAEMON_NAME="amqp2engines_test"

init 	= cinit()
logger 	= init.getLogger(DAEMON_NAME, level="INFO")
#logger 	= init.getLogger(DAEMON_NAME)
handler = init.getHandler(logger)

def main():
	global ready
		
	logger.info("Initialyze process")
	handler.run()

	# Engine_cleaner
	import perfstore2
	engine = perfstore2.engine(		next_amqp_queues=['Engine_eventstore'])

	engine.start()
	logger.info("Wait")
	handler.wait()

	engine.signal_queue.put("STOP")
	engine.join()

	logger.info("Process finished")
	
if __name__ == "__main__":
	main()
