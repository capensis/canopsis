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

from canopsis.engines.core import Engine, DROP, publish

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.event import forger, get_routingkey
from canopsis.old.mfilter import check

import json
from time import time


class engine(Engine):
    etype = 'event_filter'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        account = Account(user="root", group="root")
        self.storage = get_storage(logging_level=self.logging_level,
                                   account=account)
        self.derogations = []
        self.name = kargs['name']
        self.drop_event_count = 0
        self.pass_event_count = 0

    def pre_run(self):
        self.beat()

    def time_conditions(self, derogation):
        conditions = derogation.get('time_conditions', None)

        if not isinstance(conditions, list):
            self.logger.error(("Invalid time conditions field in '%s': %s"
                               % (derogation['_id'], conditions)))
            self.logger.debug(derogation)
            return False

        result = False

        now = time()
        for condition in conditions:
            if (condition['type'] == 'time_interval'
                    and condition['startTs']
                    and condition['stopTs']):
                always = condition.get('always', False)

                if always:
                    self.logger.debug(" + 'time_interval' is 'always'")
                    result = True

                elif (now >= condition['startTs']
                      and now < condition['stopTs']):
                    self.logger.debug(" + 'time_interval' Match")
                    result = True

        return result

    def a_override(self, event, action):
        """Override a field from event or add a new one if it does not have
        one.
        """

        afield = action.get('field', None)
        avalue = action.get('value', None)

        # This mus be a hard check because value can be a boolean or a null integer
        if afield is not None and avalue is not None:
            if afield in event and isinstance(event[afield], list):
                event[afield].append(avalue)
            else:
                event[afield] = avalue
            self.logger.debug(
                u"    + {}: Override: '{}' -> '{}'".format(
                    event['rk'], afield, avalue))
            return True

        else:
            self.logger.error(
                u"Action malformed (needs 'field' and 'value'): {}".format(
                    action))
            return False

    def a_remove(self, event, action):
        """Remove an event from a field in event or the whole field if no
        element is specified.
        """

        akey = action.get('key', None)
        aelement = action.get('element', None)
        del_met = action.get('met', 0)

        if akey:
            if aelement:
                if del_met:
                    for i, met in enumerate(event[akey]):
                        if met['name'] == aelement:
                            del event[akey][i]
                            break
                elif isinstance(event[akey], dict):
                    del event[akey][aelement]
                elif isinstance(event[akey], list):
                    del event[akey][event[akey].index(aelement)]

                self.logger.debug(u"    + {}: Removed: '{}' from '{}'".format(
                    event['rk'],
                    aelement,
                    akey))

            else:
                del event[akey]
                self.logger.debug(u"    + {}: Removed: '{}'".format(
                    event['rk'],
                    akey))

            return True

        else:
            self.logger.error(
                u"Action malformed (needs 'key' and/or 'element'): {}".format(
                    action))
            return False

    def a_modify(self, event, action, name):
        """
        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        derogated = False
        atype = action.get('type')
        actionMap = {
            'override': self.a_override,
            'remove': self.a_remove
        }

        if atype in actionMap:
            derogated = actionMap[atype](event, action)

        else:
            self.logger.warning(u"Unknown action '{}'".format(atype))

        # If the event was derogated, fill some informations
        if derogated:
            self.logger.debug(u"Event changed by rule '{}'".format(name))

        return None

    def a_drop(self, event, action, name):
        """ Drop the event.

        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        self.logger.debug(u"Event dropped by rule '{}'".format(name))
        self.drop_event_count += 1

        return DROP

    def a_pass(self, event, action, name):
        """Pass the event to the next queue.

        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        self.logger.debug(u"Event passed by rule '{}'".format(name))
        self.pass_event_count += 1

        return event

    def a_route(self, event, action, name):
        """
        Change the route to which an event will be sent
        Args:
            event: map of the event to be modified
            action: map of type action
            name: of the rule
        Returns:
            ``None``
        """

        if "route" in action:
            self.next_amqp_queues = [action["route"]]
            self.logger.debug(u"Event re-routed by rule '{}'".format(name))
        else:
            self.logger.error(
                u"Action malformed (needs 'route'): {}".format(action))

        return None

    def a_exec_job(self, event, action, name):
        records = self.storage.find(
            {'crecord_type': 'job', '_id': action['job'] }
        )
        for record in records:
            job = record.dump()
            job['context'] = event
            publish(publisher=self.amqp, event=job, rk='Engine_scheduler', exchange='amq.direct')
            #publish(publisher=self.amqp, event=job, rk='Engine_scheduler')
        return True

    def apply_actions(self, event, actions):
        pass_event = False
        actionMap = {'drop': self.a_drop,
                     'pass': self.a_pass,
                     'override': self.a_modify,
                     'remove': self.a_modify,
					 'execjob': self.a_exec_job,
                     'route': self.a_route
					}

        for name, action in actions:
            if (action['type'] in actionMap):
                ret = actionMap[action['type'].lower()](event, action, name)
                if ret:
                    pass_event = True
            else:
                self.logger.warning(u"Unknown action '{}'".format(action))

        return pass_event

    def work(self, event, *xargs, **kwargs):

        rk = get_routingkey(event)
        default_action = self.configuration.get('default_action', 'pass')

        # list of actions supported

        rules = self.configuration.get('rules', [])
        to_apply = []

        self.logger.debug(u'event {}'.format(event))

        # When list configuration then check black and
        # white lists depending on json configuration
        for filterItem in rules:
            actions = filterItem.get('actions')
            name = filterItem.get('name', 'no_name')

            self.logger.debug(u'rule {}'.format(filterItem))
            self.logger.debug(u'filter is {}'.format(filterItem['mfilter']))
            # Try filter rules on current event
            if filterItem['mfilter'] and check(filterItem['mfilter'], event):

                self.logger.debug(
                    u'Event: {}, filter matches'.format(event['rk'])
                )

                for action in actions:
                    if action['type'].lower() == 'drop':
                        self.apply_actions(event, to_apply)
                        return self.a_drop(event, None, name)
                    to_apply.append((name, action))

                if filterItem.get('break', 0):
                    self.logger.debug(
                        u' + Filter {} broke the next filters processing'
                        .format(
                            filterItem.get('name', 'filter')
                        )
                    )
                    break

        if len(to_apply):
            if self.apply_actions(event, to_apply):
                self.logger.debug(
                    u'Event before sent to next engine: %s' % event
                )
                event['rk'] = event['_id'] = get_routingkey(event)
                return event

        # No rules matched
        if default_action == 'drop':
            self.logger.debug("Event '%s' dropped by default action" % (rk))
            self.drop_event_count += 1
            return DROP

        self.logger.debug("Event '%s' passed by default action" % (rk))
        self.pass_event_count += 1

        self.logger.debug(u'Event before sent to next engine: %s' % event)
        event['rk'] = event['_id'] = get_routingkey(event)
        return event

    def beat(self, *args, **kargs):
        """ Configuration reload for realtime ui changes handling """

        self.derogations = []
        self.configuration = {
            'rules': [],
            'default_action': self.find_default_action()
        }

        self.logger.debug('Reload configuration rules')
        records = self.storage.find(
            {'crecord_type': 'filter', 'enable': True},
            sort='priority'
        )

        for record in records:

            record_dump = record.dump()
            self.set_loaded(record_dump)

            try:
                record_dump["mfilter"] = json.loads(record_dump["mfilter"])
            except Exception:
                self.logger.info('Invalid mfilter {}, filter {}'.format(
                    record_dump['mfilter'],
                    record_dump['name'],

                ))

            self.logger.debug('Loading record_dump:')
            self.logger.debug(record_dump)
            self.configuration['rules'].append(record_dump)

        self.logger.info(
            'Loaded {} rules'.format(len(self.configuration['rules']))
        )
        self.send_stat_event()

    def set_loaded(self, record):

        if 'run_once' in record and not record['run_once']:
            self.storage.update(record['_id'], {'run_once': True})
            self.logger.info(
                'record {} has been run once'.format(record['_id'])
            )

    def send_stat_event(self):
        """ Send AMQP Event for drop and pass metrics """

        message_dropped = '{} event dropped since {}'.format(
            self.drop_event_count,
            self.beat_interval
        )
        message_passed = '{} event passed since {}'.format(
            self.pass_event_count,
            self.beat_interval
        )
        event = forger(
            connector='Engine',
            connector_name='engine',
            event_type='check',
            source_type='resource',
            resource=self.amqp_queue + '_data',
            state=0,
            state_type=1,
            output=message_dropped,
            perf_data_array=[
                {'metric': 'pass_event',
                 'value': self.pass_event_count,
                 'type': 'GAUGE'},
                {'metric': 'drop_event',
                 'value': self.drop_event_count,
                 'type': 'GAUGE'}
            ]
        )

        self.logger.debug(message_dropped)
        self.logger.debug(message_passed)
        publish(publisher=self.amqp, event=event)
        self.drop_event_count = 0
        self.pass_event_count = 0

    def find_default_action(self):
        """Find the default action stored and returns it, else assume it
        default action is pass.
        """

        records = self.storage.find({'crecord_type': 'defaultrule'})
        if records:
            return records[0].dump()["action"]

        self.logger.debug(
            "No default action found. Assuming default action is pass"
        )
        return 'pass'


