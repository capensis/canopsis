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

from logging import getLogger


class Event(object):
    """
    Manage event content

    An event contains information and require a type and source_type.
    """
    SOURCE_TYPE = 'source_type'
    EVENT_TYPE = 'event_type'
    COMPONENT = 'component'
    CONNECTOR = 'connector'
    CONNECTOR_NAME = 'connector_name'


logger = getLogger('event')


def get_routingkey(event):
    """
    Build the routing key from an event.

    If the key 'resource' is present and != '', 'source_type' is forced to
    'resource', otherwise 'component'.

    This function mutates the 'source_type' field if necessary.

    :raise KeyError: on missing required info
    """
    event[Event.SOURCE_TYPE] = Event.COMPONENT
    if event.get(Event.RESOURCE, ''):
        logger.info(u"Event {} has changed source_type from {} to {}".format(event,
                    event.get(Event.SOURCE_TYPE, ''), Event.RESOURCE))
        event[Event.SOURCE_TYPE] = Event.RESOURCE

    rk = u"{}.{}.{}.{}.{}".format(
        event[Event.CONNECTOR],
        event[Event.CONNECTOR_NAME],
        event[Event.EVENT_TYPE],
        event[Event.SOURCE_TYPE],
        event[Event.COMPONENT]
    )

    if event.get(Event.RESOURCE, ''):
        rk = u"{}.{}".format(rk, event[Event.RESOURCE])

    return rk
