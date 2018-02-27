from __future__ import unicode_literals

import json

from six import string_types

def type_check(name, value, type_):
    if not isinstance(value, type_):
        raise TypeError('not a {}: {}'.format(
            type_, type(value)
        ))

    return value

class WatcherModel(object):

    def __init__(
        self, id_, entity_filter, display_name,
        description=None, enable=True, downtimes_as_ok=True, output_tpl='',
        state_when_all_ack='worststate', state_algorithm=None
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
            '_id': self.id_,
            'mfilter': json.dumps(self.entity_filter),
            'display_name': self.display_name,
            'description': self.description,
            'enable': self.enable,
            'downtimes_as_ok': self.downtimes_as_ok,
            'output_tpl': self.output_tpl,
            'state_when_all_ack': self.state_when_all_ack,
            'state_algorithm': self.state_algorithm,
        }