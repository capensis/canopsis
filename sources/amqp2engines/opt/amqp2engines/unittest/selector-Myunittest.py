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

import 	unittest
import 	sys, os
import 	logging
from 	cstorage 	import get_storage
from 	caccount 	import caccount
from 	cselector 	import cselector
import amqp

sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))
import selector


class KnownValues(unittest.TestCase):
	def setUp(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.engine = selector.engine()#logging_level=logging.DEBUG
		self.engine.storage = self.storage

	def test_01_Init(self):
		self.engine.pre_run()

		selectorTest = cselector(self.storage, name='selectorTest')
		selectorTest.mfilter = {'test_key': 'value'}
		selectorTest.load(selectorTest.dump())

		self.engine.selectors = [selectorTest]

		self.engine.work({'test_key':'not a value'})
		self.assertTrue(self.engine.selector_refresh == {})

		self.engine.selectors = [selectorTest]
		self.engine.work({'test_key':'value'})
		self.assertTrue(self.engine.selector_refresh == {'selector.account.root.selectorTest': True})

		from camqp import camqp
		self.engine.amqp = camqp(logging_level=logging.INFO, logging_name="test selector engine")

		for event_append in xrange(10):
			self.engine.beat()

		self.engine.beat()

		self.engine.post_run()


if __name__ == "__main__":
	unittest.main()



