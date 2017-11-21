#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Event standardization.
"""

from __future__ import unicode_literals

from json import dumps
from six import string_types
from time import time

DEFAULT_EVENT_TYPE = 'check'
DEFAULT_OUTPUT = ''
DEFAULT_STATE = 0
DEFAULT_STATE_TYPE = 1


class EventUTF8Error(Exception):
    pass


class Event(object):

    """
    Event representation.
    """

    _present_keys = {
        'connector',
        'connector_name',
        'component',
        'source_type',
        'event_type',
        'state',
        'state_type',
        'output'
    }

    def __init__(
            self,
            connector,
            connector_name,
            component,
            source_type,
            event_type=DEFAULT_EVENT_TYPE,
            state=DEFAULT_STATE,
            state_type=DEFAULT_STATE_TYPE,
            output=DEFAULT_OUTPUT,
            **kwargs
    ):
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.source_type = source_type
        self.event_type = event_type
        self.state = state
        self.state_type = state_type
        self.output = output

        if 'timestamp' not in kwargs:
            kwargs['timestamp'] = int(time())

        for key, val in kwargs.items():
            self._present_keys.add(key)
            setattr(self, key, val)

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
                and ('perfdata' not in self._present_keys
                     or 'perf_data_array' not in self._present_keys)):
            return False

        return True

    def ensure_utf8_format(self):
        """
        Force string attributes to be encoded in UTF-8.

        :raises: EventUTF8Error
        """
        for key in self._present_keys:
            attr = getattr(self, key)
            if isinstance(attr, string_types):
                try:
                    setattr(self, key, attr.encode('utf-8'))
                except (UnicodeEncodeError, UnicodeDecodeError):
                    raise EventUTF8Error

    def to_json(self):
        """
        Convert the Event to a json format.

        :rtype: str
        """
        try:
            return dumps({k: getattr(self, k) for k in self._present_keys})
        except ValueError:
            pass

    def __repr__(self):
        return '<Event {} {} {}>'.format(self.connector,
                                         self.connector_name,
                                         self.component)
