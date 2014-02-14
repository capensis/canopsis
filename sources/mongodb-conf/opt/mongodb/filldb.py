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

import sys
from ctools import dynmodloads
from cinit import cinit
import clogging
init = cinit()

if len(sys.argv) != 2:
	print "Usage: %s [init|update]" % sys.argv[0]
	sys.exit(1)

action = sys.argv[1].lower()

if action != "update" and action != "init":
	print "Invalid option"
	sys.exit(1)

## Logger
logger 	= clogging.getLogger()

## Load
modules = dynmodloads("~/opt/mongodb/load.d")

for name in sorted(modules):
	module = modules[name]
	module.logger = logger
	logger.info("%s %s ..." % (action, name))
	
	if action == "update":
		module.update()
	elif action == "init":
		module.init()
