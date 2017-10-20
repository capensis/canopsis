#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from unittest import main, TestCase

from canopsis.engines.cleaner import engine as Cleaner


class BaseTest(TestCase):

    def setUp(self):
        self.engine = Cleaner()


class fakeMsg():
    """
    Fake input message for Engine.work()
    """
    delivery_info = {
        "routing_key": "a/fake/rk"
    }


class TestManager(BaseTest):

    def test_state_not_an_int(self):
        body = """
{
    "connector": "canopsis",
    "connector_name": "engine",
    "event_type": "watcher",
    "source_type": "component",
    "component": "display_name",
    "resource": "resource",
    "state": "2",
    "output": "output"
}
        """

        cleaned_event = self.engine.work(body, fakeMsg())

        self.assertEqual(cleaned_event['resource'], 'resource')
        self.assertEqual(cleaned_event['state'], 2)


if __name__ == '__main__':
    main()
