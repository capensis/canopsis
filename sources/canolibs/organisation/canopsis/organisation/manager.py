#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

from md5 import new as md5

from canopsis.configuration import Parameter
from canopsis.storage.manager import Manager


class Organisation(Manager):
    """
    Dedicated to manage accounts in Canopsis
    """

    CONF_FILE = '~/etc/organisation.conf'

    CATEGORY = 'ORGANISATION'

    ACCOUNT_STORAGE = 'account_storage'
    GROUP_STORAGE = 'group_storage'
    PROFILE_STORAGE = 'profile_storage'
    RIGHTS_STORAGE = 'rights_storage'
    DATA_STORAGE = 'data_storage'

    DATA_TYPE = 'organisation'

    LOGIN = 'login'
    PWD = 'pwd'

    ACCOUNT_ID = 'account_id'
    GROUP_ID = 'group_id'
    PROFILE_ID = 'profile_id'
    RIGHT_ID = 'right_id'
    DATA_ID = 'data_id'

    ID = 'id'

    class Error(Exception):
        """
        Exception dedicated to Organisation methods errors
        """
        pass

    def __init__(
        self,
        data_type=DATA_TYPE,
        account_storage=None, group_storage=None, profile_storage=None,
        rights_storage=None, data_storage=None,
        *args, **kwargs
    ):

        super(Organisation, self).__init__(
            data_type=data_type, *args, **kwargs)

        self.account_storage = account_storage
        self.group_storage = group_storage
        self.profile_storage = profile_storage
        self.rights_storage = rights_storage
        self.data_storage = data_storage

    @property
    def account_storage(self):
        return self._account_storage

    @account_storage.setter
    def account_storage(self, value):
        self._account_storage = self._get_property_storage(value)

    @property
    def group_storage(self):
        return self._group_storage

    @group_storage.setter
    def group_storage(self, value):
        self._group_storage = self._get_property_storage(value)

    @property
    def profile_storage(self):
        return self._profile_storage

    @profile_storage.setter
    def profile_storage(self, value):
        self._profile_storage = self._get_property_storage(value)

    @property
    def rights_storage(self):
        return self._rights_storage

    @rights_storage.setter
    def rights_storage(self, value):
        self._rights_storage = self._get_property_storage(value)

    @property
    def data_storage(self):
        return self._data_storage

    @data_storage.setter
    def data_storage(self, value):
        self._data_storage = self._get_property_storage(value)

    def get_accounts(self, account_ids=None):
        """
            Get an accounts
        """

        result = self.account_storage.get(data_id=account_ids)

        return result

    def find_accounts(self, request, limit=0, skip=0, sort=None):
        """
            Find accounts which correspond to the input request

            :type request: dict
        """

        result = self.account_storage.find_elements(
            request=request, limit=limit, skip=skip, sort=sort)

        return result

    def get_account(self, login, pwd):
        """
            Get an account from a login and a password
        """

        crypted_pwd = md5(pwd)
        request = {
            Organisation.LOGIN: pwd,
            Organisation.PWD: crypted_pwd
        }
        result = self.account_storage.find(request, limit=1)

        return result

    def update_account(self, account):
        """
            Update an account
        """

        self.account_storage.update(
            data_id=account[Organisation.ID],
            value=account)

    def remove_account(self, account_id):
        """
            Remove the account which is identified by the input account_id
        """

        self.account_storage.remove(data_id=account_id)

    def get_groups(self, group_ids):
        """
        """

        result = self.group_storage.get(group_ids)

        return result

    def find_groups(self, request, limit=0, skip=0, sort=None):

        result = self.group_storage.find(request, limit, skip, sort)

        return result

    def update_group(self, group_id, value):

        self.group_storage.update(group_id, value)

    def remove_groups(self, group_ids):

        self.group_storage.remove(group_ids)

    def get_profiles(self, profile_ids):
        """
        """

        result = self.profile_storage.get(profile_ids)

        return result

    def find_profiles(self, request, limit=0, skip=0, sort=None):

        result = self.profile_storage.find(request, limit, skip, sort)

        return result

    def update_profile(self, profile_id, value):

        self.profile_storage.update(profile_id, value)

    def remove_profiles(self, profile_ids):

        self.profile_storage.remove(profile_ids)

    def get_rights(self, right_ids):
        """
        """

        result = self.right_storage.get(right_ids)

        return result

    def find_rights(self, request, limit=0, skip=0, sort=None):

        result = self.right_storage.find(request, limit, skip, sort)

        return result

    def update_right(self, right_id, value):

        self.right_storage.update(right_id, value)

    def remove_rights(self, right_ids):

        self.right_storage.remove(right_ids)

    def get_data(self, data_ids):
        """
        """

        result = self.data_storage.get(data_ids)

        return result

    def find_data(self, request, limit=0, skip=0, sort=None):

        result = self.data_storage.find(request, limit, skip, sort)

        return result

    def update_data(self, data_id, value):

        self.data_storage.update(data_id, value)

    def remove_data(self, data_ids):

        self.data_storage.remove(data_ids)

    def _conf(self, *args, **kwargs):

        result = super(Organisation, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Organisation.CATEGORY,
            new_content=(
                Parameter(Organisation.ACCOUNT_STORAGE),
                Parameter(Organisation.GROUP_STORAGE),
                Parameter(Organisation.PROFILE_STORAGE),
                Parameter(Organisation.RIGHTS_STORAGE),
                Parameter(Organisation.DATA_STORAGE)))

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(Organisation, self)._get_conf_files(*args, **kwargs)

        result.append(Organisation.CONF_FILE)

        return result
