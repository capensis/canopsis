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


class TopologyTest(TestCase):

    def setUp(self):
        self.topology = Topology()
        self.topology.data_scope = 'test'

        # set of default topologies
        self.count = 10
        self.topologies = ({Topology.ID: i} for i in range(self.count))
        # delete all topologies
        self.topology.delete()
        # check if all topologies are gone
        topologies = self.topology.get_topologies()
        self.assertFalse(topologies)

    def test_CRUD(self):

        # start to push topologies
        for topology in self.topologies:
            self.topology.push(topology=topology)

        # ensure topology count is equal to pushed topologies
        topologies = self.topology.get_topologies()
        self.assertEqual(len(topologies), self.count)

        # delete the first topology
        self.topology.delete(ids='0')
        self.assertEqual(len(topologies), self.count - 1)

        # delete third more
        self.topology.delete(ids=('-1', '2', '3'))
        self.assertEqual(len(topologies), self.count - 4)

    def test_find(self):

        for topology in self.topologies:
            self.topology.push(topology=topology)

        # find the topology where id is 3
        topologies = self.topology.find(regex='3')

        self.assertEqual(len(topologies), 1)

    def test_topologies(self):
        topologies = self.topology.get_topologies()
        self.assertFalse(topologies)

    def tearDown(self):
        # remove topologies from database
        self.topology.delete()


if __name__ == '__main__':
    main()
