#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Event standardization.
"""

from __future__ import unicode_literals

from json import dumps
from time import time

DEFAULT_RESOURCE = None
DEFAULT_ALARM = None
DEFAULT_CONTEXT = None
DEFAULT_DISPLAY_NAME = ''
DEFAULT_EVENT_TYPE = 'check'
DEFAULT_LONG_OUTPUT = ''
DEFAULT_OUTPUT = ''
DEFAULT_PERF_DATA = ''
DEFAULT_PERF_DATA_ARRAY = []
DEFAULT_STATE = 0
DEFAULT_TIMESTAMP = None
DEFAULT_MORE = {}


class EventUTF8Error(Exception):
    pass


class MinimalisticEvent(object):

    """
    Minimalistic Event representation.
    """

    def __init__(
            self,
            connector,
            connector_name,
            component,
            source_type,
            resource=DEFAULT_RESOURCE
    ):
        """
        :param str connector: connector type
        :param str connector_name: connector identifier
        :param str component: component's name
        :param str source_type: source of the event
        :param str resource: resource name
        """
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.source_type = source_type
        if resource is not None:
            self.resource = resource

    def __repr__(self):
        return '<Event {}.{}.{}>'.format(self.connector,
                                         self.connector_name,
                                         self.component)


class Event(MinimalisticEvent):

    """
    General event representation.

    Moved parameters from old spcs: tags, hostgroups, servicegroups, state_type
    """

    STR_PARAMETERS = [
        'connector',
        'connector_name',
        'component',
        'source_type',
        'resource',
        'display_name',
        'event_type',
        'long_output',
        'output',
        'perf_data'
    ]

    def __init__(
            self,
            connector,
            connector_name,
            component,
            source_type,
            resource=DEFAULT_RESOURCE,
            alarm=DEFAULT_ALARM,
            context=DEFAULT_CONTEXT,
            display_name=DEFAULT_DISPLAY_NAME,
            event_type=DEFAULT_EVENT_TYPE,
            long_output=DEFAULT_LONG_OUTPUT,
            output=DEFAULT_OUTPUT,
            perf_data=DEFAULT_PERF_DATA,
            perf_data_array=DEFAULT_PERF_DATA_ARRAY,
            state=DEFAULT_STATE,
            timestamp=DEFAULT_TIMESTAMP,
            more=DEFAULT_MORE
    ):
        """
        :param str connector: see MinimalisticEvent
        :param str connector_name: see MinimalisticEvent
        :param str component: see MinimalisticEvent
        :param str source_type: see MinimalisticEvent
        :param str resource: see MinimalisticEvent
        :param Alarm alarm: alarm informations linked to the event
        :param Context context: context informations linked to the event
        :param str display_name: name displayed in Canopsis
        :param str event_type: type of the event [sic]
        :param str long_output: long description of the event
        :param str output: short description of the event
        :param str perf_data: formatted perfdata string
        :param list perf_data_array: array of metrics
        :param int state: state level
        :param int timestamp: timestamp of the event
        :param dict more: other informations liked to the event
        """
        super(Event, self).__init__(connector=connector,
                                    connector_name=connector_name,
                                    component=component,
                                    source_type=source_type,
                                    resource=resource)
        self.alarm = alarm
        self.context = context
        self.display_name = display_name
        self.event_type = event_type
        self.long_output = long_output
        self.output = output
        self.perf_data = perf_data
        self.perf_data_array = perf_data_array
        self.timestamp = timestamp
        self.state = state
        self.more = more

        if timestamp is None:
            self.timestamp = int(time())

    def is_valid(self):
        """
        Verify that the event is valid.

        :rtype: bool
        """
        if self.source_type not in ['component', 'resource']:
            return False

        if self.source_type == 'resource' and not hasattr(self, 'resource'):
            return False

        if (self.event_type == 'perf'
                and (not hasattr(self, 'perfdata')
                     or not hasattr(self, 'perf_data_array'))):
            return False

        return True

    def ensure_utf8_format(self):
        """
        Force string attributes to be encoded in UTF-8.

        :raises: EventUTF8Error
        """
        for key in self.STR_PARAMETERS:
            if hasattr(self, key):
                attr = getattr(self, key)
                try:
                    setattr(self, key, attr.encode('utf-8'))
                except (UnicodeEncodeError, UnicodeDecodeError):
                    raise EventUTF8Error

    def to_json(self):
        """
        Convert the Event to a json format.

        :rtype: str
        """
        dico = {
            'connector': self.connector,
            'connector_name': self.connector_name,
            'component': self.component,
            'source_type': self.source_type,
            'resource': self.resource,
            'alarm': self.alarm,
            'context': self.context,
            'display_name': self.display_name,
            'event_type': self.event_type,
            'long_output': self.long_output,
            'output': self.output,
            'perf_data': self.perf_data,
            'perf_data_array': self.perf_data_array,
            'state': self.state,
            'timestamp': self.timestamp,
            'more': self.more,
        }
        try:
            return dumps(dico)
        except ValueError:
            pass
