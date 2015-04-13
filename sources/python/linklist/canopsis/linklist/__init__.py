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

__version__ = "0.1"


from canopsis.event import Event


class Check(Event):
    """
    Manage checking event with state and status information
    """

    STATE = 'state'  #: event state field name
    STATUS = 'status'  #: event status field name
    EVENT_TYPE = 'check'  #: check event type

    OK = 0  #: ok state value
    MINOR = 1  #: minor state value
    MAJOR = 2  #: major state value
    CRITICAL = 3  #: critical state value

    def __init__(self, source, state, status, meta):

        super(Event, self).__init__(
            source=source,
            data={
                Check.STATE: state,
                Check.STATUS: status
            },
            meta=meta
        )

    @property
    def state(self):
        return self.data[Check.STATE]

    @state.setter
    def state(self, value):
        self.data[Check.STATE] = value

    @property
    def status(self):
        return self.data[Check.STATUS]

    @status.setter
    def status(self, value):
        self.data[Check.STATUS] = value
