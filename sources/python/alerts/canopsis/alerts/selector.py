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

from json import loads
from logging import getLogger

from canopsis.old.cfilter import Filter
from canopsis.common.utils import singleton_per_scope
from canopsis.context.manager import Context


class StateCalculator(object):

    OFF = 0  #: off state value
    MINOR = 1  #: minor state value
    MAJOR = 2  #: major state value
    CRITICAL = 3  #: critical state value

    STATUSES = [OFF, MINOR, MAJOR, CRITICAL]

    def __init__(self, storage, logger):
        self.storage = storage
        self.logger = logger

    def get_states(self, mfilter):
        if not mfilter:
            return {}

        self.validate_filter(mfilter)

        result = self.storage.aggregate([
            {'$match': mfilter},
            {'$project': {
                '_id': True,
                'state': True
            }},
            {'$group': {
                '_id': {
                    'state': '$state',
                },
                'count': {'$sum': 1}
            }}
        ])

        self.logger.debug(u' + states : {}'.format(result))

        states = {}
        for state in result['result']:
            key = state['_id']['state']
            states[key] = states.get(key, 0) + state['count']

        self.logger.debug(u' + result: {}'.format(states))

        return states

    @property
    def storage(self):
        return self._storage

    @storage.setter
    def storage(self, value):
        self._storage = value

    @classmethod
    def gen_worst_state(cls, states):
        for s in cls.STATUSES:
            if s in states:
                yield s

    @classmethod
    def get_worst_state(cls, states, default=0):
        """Compute worst state"""

        state = default
        _gen_worst_state = cls.gen_worst_state(states)
        for s in _gen_worst_state:
            state = s
        return state

    def get_worst_state_for_ack(self, mfilter):
        if not mfilter:
            return {}

        self.validate_filter(mfilter)

        ack_clause = {'ack.isAck': {'$ne': True}, 'ack': {'$exists': True}}

        if '$and' in mfilter:
            mfilter['$and'].append(ack_clause)
        elif '$or' in mfilter:
            mfilter = {'$and': [mfilter, ack_clause]}
        elif isinstance(mfilter, dict):
            mfilter['ack'] = {'$exists': True}
            mfilter['ack.isAck'] = {'$ne': True}

        self.logger.debug(u'Selector mfilter')

        # Computes worst state for events that are not acknowleged
        states_for_ack = self.get_states(mfilter)

        self.logger.debug(u' + states for ack: {}'.format(states_for_ack))

        wstate_for_ack = StateCalculator.get_worst_state(states_for_ack)

        return wstate_for_ack

    def get_infobagor(self, mfilter):
        _, infobagot = self.storage.get_elements(
            query={
                '$and': [{'state': 0, 'status': 3}, mfilter]
            }, with_count=True)

        if infobagot:
            self.logger.info(u'infobagot count : {}'.format(infobagot))

        return infobagot

    def validate_filter(self, mfilter):
        if not isinstance(mfilter, (dict, list, tuple)):
            raise TypeError("mfilter must be a dict, list or tuple")


class Selector(object):

    TEMPLATE_REPLACE = {
        StateCalculator.OFF: '[OFF]',
        StateCalculator.MINOR: '[MINOR]',
        StateCalculator.MAJOR: '[MAJOR]',
        StateCalculator.CRITICAL: '[CRITICAL]',
    }

    DEFAULT_TEMPLATE = (
        'Off: [OFF], Minor: [MINOR], Major: [MAJOR]' +
        ', Critical: [CRITICAL]'
    )

    def __init__(
            self,
            storage,
            record=None,
            logging_level=None,
            logger=None
    ):
        self.storage = storage

        self.dostate = True
        self.data = {}
        self.mfilter = {}
        self.include_ids = []
        self.exclude_ids = []
        self.rk = None

        if logger is None:
            self.logger = getLogger('Selector')
        else:
            self.logger = logger

        if logging_level:
            self.logger.setLevel(logging_level)

        self.context = singleton_per_scope(Context)

        # Canopsis filter management for mongo
        self.cfilter = Filter()

        self.load(record.dump())

    @property
    def storage(self):
        return self._storage

    @storage.setter
    def storage(self, value):
        self._storage = value

    def load(self, dump):

        self.logger.debug(u'Loading selector from record')

        mfilter = dump.get('mfilter', '{}')

        if type(mfilter) == dict:
            self.mfilter = mfilter
        elif mfilter is None:
            self.mfilter = {}
        else:
            try:
                self.mfilter = loads(mfilter)
            except Exception as e:
                self.logger.warning(
                    u'invalid mfilter for selector {} : {} '.format(
                        dump.get('display_name', 'noname'), e))
                self.mfilter = {}

        self._id = dump.get('_id')
        self.dostate = dump.get('dostate')
        self.display_name = dump.get('display_name', 'noname')
        self.rk = dump.get('rk', self.rk)
        self.include_ids = dump.get('include_ids', self.include_ids)
        self.exclude_ids = dump.get('exclude_ids', self.exclude_ids)

        self.data = dump

        self.output_tpl = self.get_output_tpl()

        if not self.output_tpl:
            self.output_tpl = "No output template defined"

    def get_value(self, property_name, default_value):
        """
        Allow accessing record property with set of a default value
        instead of None value
        """
        # Dealing with the old strange record system.
        if (property_name not in self.data or
                self.data[property_name] is None):
            return default_value
        else:
            return self.data[property_name]

    def get_output_tpl(self):
        return self.get_value('output_tpl', Selector.DEFAULT_TEMPLATE)

    # Build MongoDB query to find every id matching event
    def makeMfilter(self):

        cfilter = self.cfilter.make_filter(
            mfilter=self.mfilter,
            includes=self.include_ids,
            excludes=self.exclude_ids,
        )

        self.logger.debug(u'Generated cfilter is')
        self.logger.debug(u'cfilter: {}'.format(cfilter))

        return cfilter

    def get_entities(self, _filter):
        """Looking for in the context of all entities that match the filter"""
        entities = self.context[Context.CTX_STORAGE].get_elements(query=_filter)
        return entities

    def alert(self):
        # Get state information form aggregation
        self.logger.debug("get state:")

        # Build MongoDB filter
        mfilter = self.makeMfilter()

        if not mfilter:
            self.logger.debug(" + Invalid filter")
            return ({}, 0, 0, 0)

        self.logger.debug(" + selector statment agregation {}".format(mfilter))

        cstate = StateCalculator(self.storage, self.logger)

        mfilter = self.makeMfilter()

        entities = self.get_entities(mfilter)

        _filter = {'entity_id': {'$in': [e['_id'] for e in entities]}}

        states = cstate.get_states(_filter)

        state = StateCalculator.get_worst_state(states)

        wstate_for_ack = cstate.get_worst_state_for_ack(mfilter)

        # Build output
        total = 0
        total_error = 0

        for s in states:
            states[s] = int(states[s])
            total += states[s]
            if s > 0:
                total_error += states[s]

        # Computed state when all events are not ack
        computed_state = wstate_for_ack

        self.logger.debug(u' + state: {}'.format(state))

        self.logger.debug(u' + total: {}'.format(total))

        output = self.output_tpl.replace('[TOTAL]', str(total))

        # output computation
        for value in StateCalculator.gen_worst_state(states):
            output = output.replace(Selector.TEMPLATE_REPLACE[value], str(value))

        data = {
            'connector': 'canopsis',
            'connector_name': 'engine',
            'component': self.display_name,
            'source_type': 'component',
            'event_type': 'selector',
            'state': computed_state,
            'output': output,
        }

        return data

    def save(self, data):
        """Store the data in database.

        :param dict data: the data to store
        """
        self.storage.put_element(data)
