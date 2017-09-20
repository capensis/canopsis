#!/usr/bin/env python
# -*- coding: utf-8  -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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
Generic event object.
"""

from __future__ import unicode_literals
from time import time


class StorageField(object):
    """
    Field names in mongo storage.
    """
    DATA_ID = '_id'
    DATA = 'd'
    TIMESTAMP = 't'
    VALUE = 'v'


class Event(object):
    """
    An event object.
    """

    def __init__(self, connector, connector_name, component, resource=None):
        """
        :param str connector:
        :param str connector_name:
        :param str component:
        :param str resource:
        """
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.resource = resource

        # TODO: handle all extra fields

        self.timestamp = int(time())
        self.data = '{}/{}'.format(self.resource, self.component)

    def to_dict(self):
        """
        Return the event as a dict.

        :rtype: dict
        """
        return {
            'connector': self.connector,
            'connector_name': self.connector_name,
            'component': self.component,
            'resource': self.resource,
        }

    def to_mongo(self):
        """
        Return the event has inserted in mongo.

        :rtype: dict
        """
        dico = {
            StorageField.DATA: self.data,
            StorageField.TIMESTAMP: self.timestamp,
            StorageField.VALUE: self.to_dict()
        }
        if hasattr(self, 'data_id'):
            dico[StorageField.DATA_ID] = self.data_id

        return dico
