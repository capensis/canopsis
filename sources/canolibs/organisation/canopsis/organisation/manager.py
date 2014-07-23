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
