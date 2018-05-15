#!/usr/bin/env python2.7
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

"""
    Mock for auth utilities (function in auth webservice, ...)
"""

from canopsis.old.account import Account
from canopsis.auth.mock.bottle import MockSession


def mock_get_account(_id=None):
    if not _id:
        _id = 'account.anonymous'

    len_prefix = len('account.')
    user = _id[len_prefix:]

    return Account(user=user)


def mock_create_session(testcase, account):
    testcase.session['account_id'] = account._id
    testcase.session['account_user'] = account.user
    testcase.session['account_group'] = account.group
    testcase.session['account_groups'] = account.groups
    testcase.session['auth_on'] = True
    testcase.session.save()

    return testcase.session


def mock_delete_session(testcase):
    testcase.session = MockSession()


def mock_check_group_rights(account, group):
    if account.group == group:
        return True

    elif group in account.groups:
        return True

    return False


def mock_check_root(account):
    if account.user == 'root':
        return True

    else:
        return mock_check_group_rights(account, 'group.CPS_root')
