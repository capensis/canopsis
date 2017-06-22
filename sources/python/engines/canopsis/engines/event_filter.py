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

from canopsis.alerts.manager import Alerts
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope
from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.event import forger, get_routingkey
from canopsis.old.mfilter import check
from canopsis.pbehavior.manager import PBehaviorManager

from json import loads


class engine(Engine):
    etype = 'event_filter'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        account = Account(user="root", group="root")
        self.storage = get_storage(logging_level=self.logging_level,
                                   account=account)
        self.name = kargs['name']
        self.drop_event_count = 0
        self.pass_event_count = 0

    def pre_run(self):
        self.beat()

    def a_override(self, event, action):
        """Override a field from event or add a new one if it does not have
        one.
        """

        afield = action.get('field', None)
        avalue = action.get('value', None)

        # This must be a hard check because value can be a boolean or a null
        # integer
        if afield is None or avalue is None:
            self.logger.error(
                "Malformed action ('field' and 'value' required): {}".format(
                    action
                )
            )
            return False

        if afield not in event:
            self.logger.debug("Overriding: '{}' -> '{}'".format(
                afield, avalue))
            event[afield] = avalue
            return True

        # afield is in event
        if not isinstance(avalue, list):
            if isinstance(event[afield], list):
                self.logger.debug("Appending: '{}' to '{}'".format(
                    avalue, afield))
                event[afield].append(avalue)

            else:
                self.logger.debug("Overriding: '{}' -> '{}'".format(
                    afield, avalue))
                event[afield] = avalue

            return True

        else:
            # operation field is supported only for list values
            op = action.get('operation', 'append')

            if op == 'override':
                self.logger.debug("Overriding: '{}' -> '{}'".format(
                    afield, avalue))
                event[afield] = avalue
                return True

            elif op == 'append':
                self.logger.debug("Appending: '{}' to '{}'".format(
                    avalue, afield))

                if isinstance(event[afield], list):
                    event[afield] += avalue
                else:
                    event[afield] = [event[afield]] + avalue

                return True

            else:
                self.logger.error(
                    "Operation '{}' unsupported (action '{}')".format(
                        op, action
                    )
                )
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
            {'crecord_type': 'job', '_id': action['job']}
        )
        for record in records:
            job = record.dump()
            job['context'] = event
            publish(
                publisher=self.amqp,
                event=job,
                rk='Engine_scheduler',
                exchange='amq.direct'
            )
            # publish(publisher=self.amqp, event=job, rk='Engine_scheduler')
        return True

    def a_snooze(self, event, action, name):
        """
        Snooze event checks

        :param dict event: event to be snoozed
        :param dict action: action
        :param str name: name of the rule

        :returns: True if a snooze has been sent, False otherwise
        :rtype: boolean
        """
        # Only check events can trigger an auto-snooze
        if event['event_type'] != 'check':
            return False

        # A check OK cannot trigger an auto-snooze
        if event['state'] == 0:
            return False

        # Alerts manager caching
        if not hasattr(self, 'am'):
            self.am = Alerts()

        # Context manager caching
        if not hasattr(self, 'cm'):
            self.cm = ContextGraph()

        entity_id = self.cm.get_id(event)

        current_alarm = self.am.get_current_alarm(entity_id)
        if current_alarm is None:
            snooze = {
                'connector': event.get('connector', ''),
                'connector_name': event.get('connector_name', ''),
                'source_type': event.get('source_type', ''),
                'component': event.get('component', ''),
                'event_type': 'snooze',
                'duration': action['duration'],
                'author': 'event_filter',
                'output': 'Auto snooze generated by rule "{}"'.format(name),
            }

            if 'resource' in event:
                snooze['resource'] = event['resource']

            publish(event=snooze, publisher=self.amqp)

            return True

        return False

    def a_baseline(self, event, actions, name):
        """a_baseline

        :param event:
        :param action: baseline conf in event filter
        :param name:
        """
        event['baseline_name'] = actions['baseline_name']
        event['check_frequency'] = actions['check_frequency']

        publish(event=event, publisher=self.amqp,
                rk='Engine_baseline', exchange='amq.direct')

    def apply_actions(self, event, actions):
        pass_event = False
        actionMap = {
            'drop': self.a_drop,
            'pass': self.a_pass,
            'override': self.a_modify,
            'remove': self.a_modify,
            'execjob': self.a_exec_job,
            'route': self.a_route,
            'snooze': self.a_snooze,
            'baseline': self.a_baseline
        }

        for name, action in actions:
            if action['type'] in actionMap:
                ret = actionMap[action['type'].lower()](event, action, name)
                if ret:
                    pass_event = True
            else:
                self.logger.warning(u"Unknown action '{}'".format(action))

        return pass_event

    def work(self, event, *xargs, **kwargs):

        rk = get_routingkey(event)
        default_action = self.configuration.get('default_action', 'pass')

        # list of supported actions

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
                    u'Event: {}, filter matches'.format(event.get('rk', event))
                )

                if 'pbehaviors' in filterItem:
                    pbehaviors = filterItem.get('pbehaviors', {})
                    list_in = pbehaviors.get('in', [])
                    list_out = pbehaviors.get('out', [])

                    if list_in or list_out:
                        pbm = singleton_per_scope(PBehaviorManager)
                        cm = singleton_per_scope(ContextGraph)
                        entity = cm.get_entity(event)
                        entity_id = cm.get_entity_id(entity)

                        result = pbm.check_pbehaviors(
                            entity_id, list_in, list_out
                        )

                        if not result:
                            break

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

        self.configuration = {
            'rules': [],
            'default_action': self.find_default_action()
        }

        self.logger.debug(u'Reload configuration rules')
        records = self.storage.find(
            {'crecord_type': 'filter', 'enable': True},
            sort='priority'
        )

        for record in records:

            record_dump = record.dump()
            self.set_loaded(record_dump)

            try:
                record_dump["mfilter"] = loads(record_dump["mfilter"])
            except Exception:
                self.logger.info(u'Invalid mfilter {}, filter {}'.format(
                    record_dump['mfilter'],
                    record_dump['name'],

                ))

            self.logger.debug(u'Loading record_dump:')
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
