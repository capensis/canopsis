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

from canopsis.topology.elements import Node
from canopsis.topology.rule.condition import new_state, condition, _all


class NewStateTest(TestCase):
    """
    Test new_state function.
    """

    def setUp(self):

        self.node = Node()

    def test_new_state_from_event(self):
        """
        Test if state is different from an event.
        """

        result = new_state(node=self.node, event={'state': 1})

        self.assertTrue(result)

    def test_same_state_from_event(self):
        """
        Test if state is the same than the one from an event.
        """

        result = new_state(node=self.node, event={'state': 0})

        self.assertFalse(result)

    def test_new_state_from_state(self):
        """
        Test if state is different than the input state.
        """

        result = new_state(node=self.node, event={}, state=1)

        self.assertTrue(result)

    def test_same_state_from_state(self):
        """
        Test if state is the same than the input state.
        """

        result = new_state(node=self.node, event={}, state=0)

        self.assertFalse(result)


class ConditionTest(TestCase):
    pass


class AllTest(TestCase):
    pass


if __name__ == '__main__':
    main()
