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

from canopsis.check import Check
from canopsis.check.manager import CheckManager


class CheckManagerTest(TestCase):
    """
    Base class for all check manager tests.
    """

    def setUp(self):
        """
        initialize a manager.
        """

        self.manager = CheckManager(data_scope='test_check')

    def tearDown(self):

        self.manager.del_state()


class StateTest(CheckManagerTest):
    """
    Test manager.state method.
    """

    def test_not_existing_entity(self):
        """
        Test to get state from an entity which does not exist.
        """

        entity_id = 'test'

        state = self.manager.state(ids=entity_id)

        self.assertIsNone(state)

    def test_existing_entity(self):
        """
        Test to get a state from an old value.
        """
        entity_id = 'test'

        state = Check.OK

        new_state = self.manager.state(ids=entity_id, state=state)
        self.assertEqual(new_state, state)

        state += 1  # increment state
        # change state of entity and check the result value
        new_state = self.manager.state(ids=entity_id, state=state)
        self.assertEqual(state, new_state)
        # check if DB state equals new updated state
        new_state = self.manager.state(ids=entity_id)
        self.assertEqual(state, new_state)

    def test_hard_criticity(self):
        """
        Test hard criticity.
        """

        entity_id = 'test'

        hard = CheckManager.HARD
        ok = Check.OK
        warning = Check.WARNING

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, None)

        state = self.manager.state(ids=entity_id, criticity=hard, state=ok)
        self.assertEqual(state, ok)

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, ok)

        state = self.manager.state(
            ids=entity_id, criticity=hard, state=warning
        )
        self.assertEqual(state, warning)

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, warning)

    def test_soft_criticity(self):
        """
        Test soft criticity.
        """

        entity_id = 'test'
        soft = CheckManager.SOFT
        soft_count = CheckManager.CRITICITY_COUNT[soft]
        ok = Check.OK
        warning = Check.WARNING
        unknown = Check.UNKNOWN

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, None)

        state = self.manager.state(ids=entity_id, criticity=soft, state=ok)
        self.assertEqual(state, Check.UNKNOWN)

        for i in range(1, soft_count - 1):
            state = self.manager.state(
                ids=entity_id, criticity=soft, state=ok
            )
            self.assertEqual(state, unknown)
            state = self.manager.state(ids=entity_id)
            self.assertEqual(state, unknown)

        state = self.manager.state(ids=entity_id, criticity=soft, state=ok)
        self.assertEqual(state, ok)
        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, ok)

        for i in range(soft_count - 1):
            state = self.manager.state(
                ids=entity_id, criticity=soft, state=warning
            )
            self.assertEqual(state, ok)
            state = self.manager.state(ids=entity_id)
            self.assertEqual(state, ok)

        state = self.manager.state(
            ids=entity_id, criticity=soft, state=warning
        )
        self.assertEqual(state, warning)
        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, warning)


if __name__ == '__main__':
    main()
