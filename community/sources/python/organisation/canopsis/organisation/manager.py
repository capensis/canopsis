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

from md5 import new as md5

from canopsis.configuration.model import Parameter
from canopsis.middleware.registry import MiddlewareRegistry


class Organisation(MiddlewareRegistry):
    """
    Dedicated to manage users in Canopsis
    """

    CONF_FILE = 'organisation/organisation.conf'

    CATEGORY = 'ORGANISATION'

    USER_STORAGE = 'user_storage'
    RIGHT_STORAGE = 'right_storage'
    PROFILE_STORAGE = 'profile_storage'

    DATA_TYPE = 'organisation'

    LOGIN = 'login'
    PWD = 'pwd'

    USER_ID = 'user_id'
    RIGHT_ID = 'right_id'
    PROFILE_ID = 'profile_id'

    ID = 'id'

    CHECKSUM = 'checksum'

    CONTEXT = 'context'

    class Error(Exception):
        """
        Exception dedicated to Organisation methods errors
        """

    def __init__(
        self,
        data_type=DATA_TYPE,
        user_storage=None, right_storage=None, profile_storage=None,
        group_storage=None,
        *args, **kwargs
    ):

        super(Organisation, self).__init__(
            data_type=data_type, *args, **kwargs)

        self.user_storage = self.get_storage()
        self.right_storage = self.get_storage()
        self.profile_storage = self.get_storage()
        self.group_storage = self.get_typed_storage()

    @property
    def user_storage(self):

        return self._user_storage

    @user_storage.setter
    def user_storage(self, value):

        self._user_storage = self._get_property_storage(value)

    @property
    def right_storage(self):

        return self._group_storage

    @right_storage.setter
    def right_storage(self, value):

        self._group_storage = self._get_property_storage(value)

    @property
    def profile_storage(self):

        return self._profile_storage

    @profile_storage.setter
    def profile_storage(self, value):

        self._profile_storage = self._get_property_storage(value)

    def get_users(self, user_ids=None):
        """
        Get an users
        """

        result = self.user_storage.get(data_id=user_ids)

        return result

    def find_users(self, request, limit=0, skip=0, sort=None):
        """
        Find users which correspond to the input request

        :type request: dict
        """

        result = self.user_storage.find_elements(
            request=request, limit=limit, skip=skip, sort=sort)

        return result

    def get_user(self, login, pwd):
        """
        Get an user from a login and a password.
        """

        crypted_pwd = md5(pwd)
        request = {
            Organisation.LOGIN: pwd,
            Organisation.PWD: crypted_pwd
        }
        result = self.user_storage.find(request, limit=1)

        return result

    def update_user(self, user):
        """
        Update an user.
        """

        self.user_storage.update(data_id=user[Organisation.ID], value=user)

    def remove_users(self, user_ids):
        """
        Remove users identified by the input user_ids
        """

        self.user_storage.remove(data_ids=user_ids)

    def get_groups(self, group_ids):
        """
        """

        result = self.group_storage.get(data_ids=group_ids)

        return result

    def find_groups(self, request, _type=None, limit=0, skip=0, sort=None):

        result = self.group_storage.find(request, limit, skip, sort)

        return result

    def update_group(self, group_id, value):

        self.group_storage.update(group_id, value)

    def remove_groups(self, group_ids=None, _type=None):

        self.group_storage.remove(data_ids=group_ids, _type=_type)

    def get_profiles(self, profile_ids):
        """
        """

        result = self.profile_storage.get(data_ids=profile_ids)

        return result

    def find_profiles(self, request, limit=0, skip=0, sort=None):

        result = self.profile_storage.find(request, limit, skip, sort)

        return result

    def update_profile(self, profile_id, value):

        self.profile_storage.update(data_id=profile_id, value=value)

    def remove_profiles(self, profile_ids):

        self.profile_storage.remove(data_ids=profile_ids)

    def get_rights(self, right_ids):
        """
        """

        result = self.right_storage.get(data_ids=right_ids)

        return result

    def find_rights(self, request, limit=0, skip=0, sort=None):

        result = self.right_storage.find(request, limit, skip, sort)

        return result

    def update_right(self, right_id, value):

        self.right_storage.update(data_id=right_id, value=value)

    def remove_rights(self, right_ids):

        self.right_storage.remove(data_ids=right_ids)

    def check_rights(self, user, permission):
        """
        Check whatever or not if user rights checked input permission
        """

    def _conf(self, *args, **kwargs):

        result = super(Organisation, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Organisation.CATEGORY,
            new_content=(
                Parameter(Organisation.USER_STORAGE),
                Parameter(Organisation.RIGHT_STORAGE),
                Parameter(Organisation.PROFILE_STORAGE)))

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(Organisation, self)._get_conf_files(*args, **kwargs)

        result.append(Organisation.CONF_FILE)

        return result
