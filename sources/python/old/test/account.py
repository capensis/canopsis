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

import unittest

from canopsis.old.account import Account, caccount_get
from canopsis.old.group import Group
from canopsis.old.storage import Storage

STORAGE = None
ACCOUNT = None
GROUP = None


class KnownValues(unittest.TestCase):
    def setUp(self):
        pass

    def test_01_Init(self):
        global ACCOUNT
        ACCOUNT = Account(user="wpain", lastname="Pain", firstname="William", mail="wpain@capensis.fr", group="capensis", groups=['group.titi', 'group.tata'])
        global GROUP

        user_account = Account(user="william", group="capensis")
        user_account.cat()

    def test_02_Cat(self):
        ACCOUNT.cat()

    def test_03_Passwd(self):
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

        ACCOUNT.cat()

    def test_04_authkey(self):
        authkey = ACCOUNT.get_authkey()
        if not authkey:
            raise Exception('Invalid authkey ... (%s)' % authkey)

    def test_05_Store(self):
        STORAGE.put(ACCOUNT)
        ACCOUNT.cat()
    """
    def test_05_GetAll(self):
        account = Account(user="ojan", lastname="Jan", firstname="Olivier", mail="ojan@capensis.fr", group="capensis")
        STORAGE.put(account)

        accounts = caccount_getall(STORAGE)

        if len(accounts) != 2:
            raise Exception('caccount_getall dont work ...')
    """
    def test_06_Edit(self):
        ACCOUNT.chgrp('toto')
        ACCOUNT.cat()
        STORAGE.put(ACCOUNT)

    def test_07_CheckGet(self):
        record = STORAGE.get("account.wpain")
        record.cat()
        account = Account(record)
        account.cat()

        if account.user != 'wpain':
            raise Exception('account.user: Corruption in load ...')

        if account.group != 'group.toto':
            raise Exception('account.group: Corruption in load ...')

        if account.groups != ['group.titi', 'group.tata']:
            raise Exception('account.groups: Corruption in load ...')

    def test_08_CheckEdit(self):
        account = caccount_get(STORAGE, "wpain")

        if account.group != 'group.toto':
            raise Exception('Impossible to edit account in DB ...')

    def test_09_Remove(self):
        ## Anonymous cant remove account
        self.assertRaises(ValueError, STORAGE.remove, ACCOUNT, Account())

        ## But root can ;)
        STORAGE.remove(ACCOUNT)
    """
    def test_10_check_addgroup_removegroup(self):
        GROUP = Group(name='mygroup')
        ACCOUNT.add_in_groups(GROUP)

        if GROUP._id not in ACCOUNT.groups:
            raise Exception('Error while add_in_groups, group not added')
        if ACCOUNT._id not in GROUP.account_ids:
            raise Exception('Error while add_in_groups, account not added to group')

        ACCOUNT.remove_from_groups(GROUP)

        if GROUP._id in ACCOUNT.groups:
            raise Exception('Error while remove_from_groups, group not removed')
        if ACCOUNT._id in GROUP.account_ids:
            raise Exception('Error while remove_from_groups, group not removed from account')
    """
    def test_11_check_group_func_autosav(self):
        account = Account(user='test', lastname='testify', storage=STORAGE)
        group = Group(name='Mgroup')

        STORAGE.put(account)
        STORAGE.put(group)

        account.add_in_groups(group._id)

        bdd_account = Account(STORAGE.get(account._id))
        bdd_group = Group(STORAGE.get(group._id))

        if group._id not in bdd_account.groups:
            raise Exception('Group corruption while stock in bdd after add in group')
        if account._id not in bdd_group.account_ids:
            raise Exception('Group corruption while stock in bdd after add in group')
        '''
        account.remove_from_groups(group._id)

        bdd_account = Account(STORAGE.get(bdd_account._id))
        bdd_group = Group(STORAGE.get(bdd_group._id))

        if group._id in bdd_account.groups:
            raise Exception('Group corruption while stock in bdd after remove from group')
        if account._id in bdd_group.account_ids:
            raise Exception('Group corruption while stock in bdd after remove from group')
        '''
    def test_99_DropNamespace(self):
        STORAGE.drop_namespace('unittest')

if __name__ == "__main__":
    STORAGE = Storage(Account(user="root", group="root"), namespace='unittest')
    unittest.main(verbosity=2)
