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


class RightEntity(object):
    """
    Right entity
    """

    def __init__(self, _id, name, *rights):

        super(RightEntity, self).__init__()

        self.name = name
        self.id = _id
        self.rights = rights


class Right(object):
    """
    In charge of managing permission with an input permission value
    """

    NONE = 0
    READ = 1
    UPDATE = READ << 1
    DELETE = UPDATE << 1
    CREATE = DELETE << 1

    ALL = READ | UPDATE | DELETE | CREATE

    def __init__(
        self, _id, name=None, checksum=NONE, element_id=None, context=None
    ):

        super(Right, self).__init__()

        self.id = _id
        self.name = name
        self.checksum = checksum
        self.element_id = element_id
        self.context = context

    def get_checksum(self, element_id, context):
        """
        Check if input value grants access to input permission

        :return: True iif value is not None and value >= permission
        :rtype: bool
        """

        raise NotImplementedError()


class RightComposite(Right, RightEntity):

    def __init__(self, rights=None, **kwargs):

        super(RightComposite, self).__init__(**kwargs)

        self.rights = rights
