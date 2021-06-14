from __future__ import unicode_literals

import json

from six import string_types


def type_check(name, value, type_):
    if not isinstance(value, type_):
        raise TypeError('not a {}: {}'.format(
            type_, type(value)
        ))

    return value

WORSTSTATE = 'worststate'


class WatcherModel(object):
    """
    Model to ease making Watchers.
    """

    _ID = '_id'
    MFILTER = 'mfilter'
    DISPLAY_NAME = 'display_name'
    DESCRIPTION = 'description'
    ENABLE = 'enable'
    DOWNTIMES_AS_OK = 'downtimes_as_ok'
    OUTPUT_TPL = 'output_tpl'
    STATE_WHEN_ALL_ACK = 'state_when_all_ack'
    STATE_ALGORITHM = 'state_algorithm'

    def __init__(
        self, id_, entity_filter, display_name,
        description=None, enable=True, downtimes_as_ok=True, output_tpl='',
        state_when_all_ack=WORSTSTATE, state_algorithm=None
    ):
        if description is None:
            description = display_name
        if state_algorithm is not None:
            type_check('state_algorithm', state_algorithm, string_types)

        self.id_ = type_check("id_", id_, string_types)
        self.entity_filter = type_check('entity_filter', entity_filter, dict)
        self.display_name = type_check('display_name', display_name, string_types)
        self.description = type_check('description', description, string_types)
        self.enable = type_check('enable', enable, bool)
        self.downtimes_as_ok = type_check('downtimes_as_ok', downtimes_as_ok, bool)
        self.output_tpl = type_check('output_tpl', output_tpl, string_types)
        self.state_when_all_ack = type_check('state_when_all_ack', state_when_all_ack, string_types)
        self.state_algorithm = state_algorithm

    def to_dict(self):
        return {
            self._ID: self.id_,
            self.MFILTER: json.dumps(self.entity_filter),
            self.DISPLAY_NAME: self.display_name,
            self.DESCRIPTION: self.description,
            self.ENABLE: self.enable,
            self.DOWNTIMES_AS_OK: self.downtimes_as_ok,
            self.OUTPUT_TPL: self.output_tpl,
            self.STATE_WHEN_ALL_ACK: self.state_when_all_ack,
            self.STATE_ALGORITHM: self.state_algorithm,
        }
