#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.topology.manager import Topology
from canopsis.topology.process import event_processing


class ProcessingTest(TestCase):

    def setUp(self):
        self.topology = Topology(data_scope='test')
        self.check = {
            'event_type': 'check',
            'connector': '',
            'connector_name': '',
            'component': '',
            'source_type': 'component'}

    def test_no_bound(self):
        """
        Test in case of not bound nodes
        """

        event = event_processing(event=self.check)
        self.assertEqual(self.check, event)

    def test_one_node(self):
        """
        Test in case of one bound node
        """

        raise NotImplementedError()

    def test_nexts(self):
        """
        Test next nodes
        """

        raise NotImplementedError()

    def tearDown(self):
        self.topology.delete_nodes()


if __name__ == '__main__':
    main()
