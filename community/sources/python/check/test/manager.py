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

from canopsis.check import Check
from canopsis.check.manager import CheckManager, InvalidState


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
        minor = Check.MINOR

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, None)

        state = self.manager.state(ids=entity_id, criticity=hard, state=ok)
        self.assertEqual(state, ok)

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, ok)

        state = self.manager.state(
            ids=entity_id, criticity=hard, state=minor
        )
        self.assertEqual(state, minor)

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, minor)

    def test_soft_criticity(self):
        """
        Test soft criticity.
        """

        entity_id = 'test'
        soft = CheckManager.SOFT
        soft_count = CheckManager.CRITICITY_COUNT[soft]
        ok = Check.OK
        minor = Check.MINOR
        critical = Check.CRITICAL

        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, None)

        state = self.manager.state(ids=entity_id, state=critical)
        self.assertEqual(state, Check.CRITICAL)

        state = self.manager.state(ids=entity_id, criticity=soft, state=ok)
        self.assertEqual(state, Check.CRITICAL)

        for i in range(1, soft_count - 1):
            state = self.manager.state(
                ids=entity_id, criticity=soft, state=ok
            )
            self.assertEqual(state, critical)
            state = self.manager.state(ids=entity_id)
            self.assertEqual(state, critical)

        state = self.manager.state(ids=entity_id, criticity=soft, state=ok)
        self.assertEqual(state, ok)
        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, ok)

        for i in range(soft_count - 1):
            state = self.manager.state(
                ids=entity_id, criticity=soft, state=minor
            )
            self.assertEqual(state, ok)
            state = self.manager.state(ids=entity_id)
            self.assertEqual(state, ok)

        state = self.manager.state(
            ids=entity_id, criticity=soft, state=minor
        )
        self.assertEqual(state, minor)
        state = self.manager.state(ids=entity_id)
        self.assertEqual(state, minor)

    def clean(self):
        # ensure state collection clean
        self.manager.del_state()
        states = self.manager.get_state()
        self.assertEqual(list(states), [])

##TODO4-01-2017
#    def test_put_state(self):
#        self.clean()
#
#        # Test can start here
#        # Test is : insert a state an retrieve it then check data content
#        self.manager.put_state('state_id', 1)
#        states = list(self.manager[CheckManager.CHECK_STORAGE].get_elements(
#            ids=None
#        ))
#        self.assertEqual(len(states), 1)
#        self.assertEqual(states[0]['state'], 1)
#        self.assertEqual(states[0]['_id'], 'state_id')
#
#        # Tests state change for given id
#        self.manager.put_state('state_id', 0)
#
#        states = list(self.manager[CheckManager.CHECK_STORAGE].get_elements(
#            ids=None
#        ))
#        self.assertEqual(states[0]['state'], 0)
#
#        # Test valid state values
#        for wrong_value in [-1, 4, True, 'test', object()]:
#            def test_state_raises():
#                self.manager.put_state('state_id', wrong_value)
#            self.assertRaises(InvalidState, test_state_raises)
#
#        # Tests no exception raised for valid states
#        for state in [0, 1, 2, 3]:
#            self.manager.put_state('state_id', state)
#
    def test_get_state(self):
        self.clean()

        self.manager[CheckManager.CHECK_STORAGE].put_element(
            _id='entity_id_1', element={'state': 1}
        )
        self.manager[CheckManager.CHECK_STORAGE].put_element(
            _id='entity_id_2', element={'state': 2}
        )

        # Test single element matching
        states = list(self.manager.get_state(ids=['entity_id_1']))
        self.assertEqual(len(states), 1)
        self.assertIn(states[0]['_id'], 'entity_id_1')

        states = list(self.manager.get_state())

        # Test we have 2 state saved and value are properly set
        self.assertEqual(len(states), 2)
        for state in states:
            self.assertIn(state['state'], [1, 2])
            self.assertIn(state['_id'], ['entity_id_1', 'entity_id_2'])

    def test_del_state(self):
        self.clean()

        self.manager[CheckManager.CHECK_STORAGE].put_element(
            _id='entity_id', element={'state': 1}
        )
        states = list(self.manager.get_state())
        self.assertEqual(len(states), 1)
        self.manager.del_state(['entity_id'])
        states = list(self.manager.get_state())
        self.assertEqual(len(states), 0)

if __name__ == '__main__':
    main()
