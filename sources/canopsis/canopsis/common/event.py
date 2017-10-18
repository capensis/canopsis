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


class Event(object):
    """
    An event object.
    """

    _known_attributes = ['connector', 'connector_name', 'component']

    def __init__(self, connector, connector_name, component, **kwargs):
        """
        :param str connector:
        :param str connector_name:
        :param str component:
        """
        self.connector = connector
        self.connector_name = connector_name
        self.component = component

        self.set(kwargs)

    def __setattr__(self, key, value):
        if key not in self._known_attributes:
            self._known_attributes.append(key)
        super(Event, self).__setattr__(key, value)

    def set(self, dico):
        """
        Update event values from a dict.

        :param dict dico: a dict of attributes to update
        """
        for key, val in dico.items():
            setattr(self, key, val)

    def to_dict(self):
        """
        Return the event as a dict.

        :rtype: dict
        """
        dico = {}
        for key in self._known_attributes:
            if hasattr(self, key):
                dico[key] = getattr(self, key)

        return dico
