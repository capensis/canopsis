# -*- coding: utf-8 -*-
# --------------------------------
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

__version__ = "0.1"

from canopsis.old.event import forger as eforger


class Event(object):
    """
    Manage event content

    An event contains information and require a type and source_type.
    """

    TYPE = 'event_type'  #: event type field name
    SOURCE_TYPE = 'source_type'  #: source type field name
    SOURCE = 'source'  #: source field name
    DATA = 'data'  #: data field name
    META = 'meta'  #: meta field name

    CONNECTOR = 'canopsis'  #: default connector value
    CONNECTOR_NAME = 'engine'  #: default connector name

    __slots__ = (TYPE, SOURCE, DATA, META)

    def __init__(self, source, data, meta, _type=None):

        super(Event, self).__init__()

        self.type = type(self).__name__.lower() if _type is None else _type
        self.source = source
        self.data = data
        self.meta = meta

    @classmethod
    def new_event(event_class, **old_event):
        """
        Create an Event from an old event (ficus and older version).
        """

        _type = event_class.__name__.lower()
        _type = old_event.pop(Event.EVENT_TYPE, _type)
        source = old_event.pop(Event.SOURCE)
        data = old_event.pop(Event.DATA, None)
        meta = old_event.pop(Event.META, None)

        result = Event(
            _type=_type,
            source=source,
            data=data,
            meta=meta)

        return result

    @classmethod
    def get_type(cls):
        """
        Get unique event type name
        """

        result = cls.__name__.lower()

        return result


def forger(
    event_type,
    connector=None, connector_name=None, component=None, resource=None,
    **kwargs
):
    """
    Forge an event from input parameters.
    """

    # init parameters
    if connector is None:
        connector = Event.CONNECTOR
    if connector_name is None:
        connector_name = Event.CONNECTOR_NAME
    # construct the event
    result = eforger(
        connector=connector,
        connector_name=connector_name,
        event_type=event_type,
        component=component
    )
    # try to put resource
    if resource is not None:
        result['resource'] = resource
        result['source_type'] = 'resource'
    else:
        # or specify source type is a component
        result['source_type'] = 'component'

    result.update(kwargs)

    return result
