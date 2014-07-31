#!/usr/bin/env python2.7
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

import unittest

from canopsis.old.account import Account
import canopsis.auth.base as base

import canopsis.auth.mock as mock


class TestBaseBackend(unittest.TestCase):
    def mock_canopsis(self):
        base.check_root = mock.auth.mock_check_root
        base.check_group_rights = mock.auth.mock_check_group_rights
        base.create_session = mock.auth.mock_create_session

    def setUp(self):
        self.backend = base.BaseBackend(
            ['group.CPS_testgroup_allow'],
            ['group.CPS_testgroup_disallow']
        )

        self.mock_canopsis()

    def test_install_account_ok(self):
        account = Account(
            user='canotest',
            group='group.CPS_testgroup_allow'
        )

        res = self.backend.install_account(account)

        self.assertTrue(res)

    def test_install_account_ko(self):
        account = Account(
            user='canotest',
            group='group.CPS_testgroup_disallow'
        )

        res = self.backend.install_account(account)

        self.assertFalse(res)
