from time import time
from six import string_types


class EventUTF8Error(Exception):
    pass


class EventUTF8Error(Exception):
    pass


class Event(object):

    _present_keys = {
        'connector',
        'connector_name',
        'component',
        'source_type',
        'event_type'
    }

    def __init__(
            self,
            connector,
            connector_name,
            component,
            source_type,
            event_type,
            **kwargs
    ):
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.source_type = source_type
        self.event_type = event_type

        if 'timestamp' not in kwargs:
            kwargs['timestamp'] = int(time())

        if 'state' not in kwargs:
            kwargs['state'] = 0

        if 'state_type' not in kwargs:
            kwargs['state_type'] = 1

        if 'event_type' not in kwargs:
            kwargs['event_type'] = 'check'

        if 'output' not in kwargs:
            kwargs['output'] = ''

        for key, val in kwargs.items():
            self._present_keys.add(key)
            setattr(self, key, val)

    def is_valid(self):
        if self.event_type == 'resource' and not hasattr(self, 'resource'):
            return False
        if (self.event_type == 'perf'
                and ('perfdata' not in self._present_keys
                     or 'perf_data_array' not in self._present_keys)):
            return False

        return True

    def ensure_utf8_format(self):
        for key in self._present_keys:
            attr = getattr(self, key)
            if isinstance(attr, string_types):
                try:
                    setattr(self, key, attr.encode('utf-8'))
                except UnicodeEncodeError, UnicodeDecodeError:
                    raise EventUTF8Error

    def to_json(self):
        try:
            return loads({k:getattr(self, k) for k in self._present_keys})
        except ValueError:
            pass

    def __repr__(self):
        return '<event {} {} {}>'.format(self.connector, self.connector_name, self.component) 
