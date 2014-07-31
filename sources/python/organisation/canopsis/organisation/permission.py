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


class PermissionEntity(object):
    """
    Permission entity
    """

    def __init__(self, _id, name, *permissions):

        super(PermissionEntity, self).__init__()

        self.name = name
        self.id = _id
        self.permissions = permissions


class Permission(object):
    """
    In charge of managing permission with an input permission value
    """

    READ = 1
    UPDATE = 2
    DELETE = 3
    CREATE = 4

    def __init__(self, _id, name, value, *args, **kwargs):

        super(Permission, self).__init__()

        self.id = _id
        self.name = name
        self.value = value

    @staticmethod
    def check(value, permission):
        """
        Check if input value grants access to input permission

        :return: True iif value is not None and value >= permission
        :rtype: bool
        """

        return value is not None and value >= permission
