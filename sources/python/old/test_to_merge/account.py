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

import unittest

from canopsis.old.account import Account, caccount_get
from canopsis.old.storage import Storage

STORAGE = None
ACCOUNT = None
GROUP = None


class KnownValues(unittest.TestCase):

    def test_01_Init(self):
        user_account = Account(user="william", group="capensis")

    def test_03_Passwd(self):
        ACCOUNT = Account(user="william", group="capensis")

        passwd = 'root'
        ACCOUNT.passwd(passwd)

        shadow = ACCOUNT.make_shadow(passwd)
        if not ACCOUNT.check_shadowpasswd(shadow):
            raise Exception('Invalid shadow passwd ... (%s)' % shadow )

        if not ACCOUNT.check_passwd(passwd):
            raise Exception('Invalid passwd ... (%s)' % passwd)

        cryptedKey = ACCOUNT.make_tmp_cryptedKey()
        if not ACCOUNT.check_tmp_cryptedKey(cryptedKey):
            raise Exception('Invalid cryptedKey ... (%s)' % authkey)

    def test_04_authkey(self):
        ACCOUNT = Account(user="william", group="capensis")

        authkey = ACCOUNT.get_authkey()
        if not authkey:
            raise Exception('Invalid authkey ... (%s)' % authkey)

    def test_05_Store(self):
        ACCOUNT = Account(user="william", group="capensis")

        STORAGE.put(ACCOUNT)

    def test_09_Remove(self):
        # Anonymous cant remove account
        self.assertRaises(ValueError, STORAGE.remove, ACCOUNT, Account())

        # But root can ;)
        STORAGE.remove(ACCOUNT)

    def test_99_DropNamespace(self):
        STORAGE.drop_namespace('unittest')

if __name__ == "__main__":
    STORAGE = Storage(Account(user="root", group="root"), namespace='unittest')
    unittest.main(verbosity=2)
