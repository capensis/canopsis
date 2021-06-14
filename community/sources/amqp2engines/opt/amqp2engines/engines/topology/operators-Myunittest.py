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

import unittest, sys, os
import logging

logging.basicConfig(level=logging.DEBUG,
	format='%(name)s %(levelname)s %(message)s',
)

logger = logging.getLogger("operators")

sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/topology/'))

def test_op(operator, states, options, value):
	state = operator(states=states, options=options)
	if state != value:
		raise Exception("Options: %s States: %s | Get state %s, expected: %s" % (options, states, state, value))

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_op_and(self):
		logger.info("Test AND")
		module = __import__('and')
		module.logger = logger

		options = { 'state': 0, 'then': 0, 'else': -1}

		# Cases
		test_op(module.operator, [0, 0, 0], options, 0)
		test_op(module.operator, [0, 0, 1], options, 1)
		test_op(module.operator, [0, 1, 1], options, 1)
		test_op(module.operator, [1, 1, 1], options, 1)
		test_op(module.operator, [0, 0, 2], options, 2)
		test_op(module.operator, [0, 2, 2], options, 2)
		test_op(module.operator, [2, 2, 2], options, 2)

	def test_op_or(self):
		logger.info("Test OR")
		module = __import__('or')
		module.logger = logger

		options = { 'state': 0, 'then': 0, 'else': -1}

		# Cases
		test_op(module.operator, [0, 0, 0], options, 0)
		test_op(module.operator, [0, 0, 1], options, 0)
		test_op(module.operator, [0, 1, 1], options, 0)
		test_op(module.operator, [1, 1, 1], options, 1)
		test_op(module.operator, [0, 0, 2], options, 0)
		test_op(module.operator, [0, 2, 2], options, 0)
		test_op(module.operator, [2, 2, 2], options, 2)

	def test_op_worst_state(self):
		logger.info("Test WORST STATE")
		module = __import__('worst_state')
		module.logger = logger

		options = { 'state': 0, 'then': 0, 'else': -1}

		# Cases
		test_op(module.operator, [0, 0, 0], options, 0)
		test_op(module.operator, [0, 0, 1], options, 1)
		test_op(module.operator, [0, 1, 1], options, 1)
		test_op(module.operator, [1, 1, 1], options, 1)
		test_op(module.operator, [0, 0, 2], options, 2)
		test_op(module.operator, [0, 2, 2], options, 2)
		test_op(module.operator, [2, 2, 2], options, 2)

	def test_op_cluster(self):
		logger.info("Test CLUSTER")
		module = __import__('cluster')
		module.logger = logger

		options = { 'least': 2, 'state': 0, 'then': 0, 'else': -1}

		# Cases
		test_op(module.operator, [0, 0, 0], options, 0)
		test_op(module.operator, [0, 0, 1], options, 0)
		test_op(module.operator, [0, 1, 1], options, 1)
		test_op(module.operator, [1, 1, 1], options, 1)
		test_op(module.operator, [0, 0, 2], options, 0)
		test_op(module.operator, [0, 2, 2], options, 2)
		test_op(module.operator, [2, 2, 2], options, 2)

if __name__ == "__main__":
	unittest.main(verbosity=2)