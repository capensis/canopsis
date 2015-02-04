#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

from unittest import TestCase, main

from sys import path

from os.path import expanduser

from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.engines.selector import engine

path.append(expanduser('~/opt/amqp2engines/engines/'))


class KnownValues(TestCase):
    def setUp(self):
        self.storage = get_storage(
            namespace='object', account=Account(user="root", group="root"))
        self.engine = engine()  # logging_level=logging.DEBUG
        self.engine.storage = self.storage

    def test_01_Init(self):
        self.engine.pre_run()
        """
        selectorTest = Selector(self.storage, name='selectorTest')
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
        """

if __name__ == "__main__":
    main()
