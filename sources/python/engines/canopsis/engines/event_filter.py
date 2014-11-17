# -*- coding: utf-8 -*-
#--------------------------------
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

from canopsis.engines import Engine, DROP

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.event import forger, get_routingkey
from canopsis.old.mfilter import check

from ast import literal_eval
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

    def pre_run(self):
        """
        Init vars and call beat method
        """

        self.drop_event_count = 0
        self.pass_event_count = 0
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
        """
        Override a field from event or add a new one if it does not have one
        """

        afield = action.get('field', None)
        avalue = action.get('value', None)

        if afield and avalue:
            if afield in event and isinstance(event[afield], list):
                event[afield].append(avalue)
            else:
                event[afield] = avalue
            self.logger.debug(("    + %s: Override: '%s' -> '%s'"
                           % (event['rk'], afield, avalue)))
            return True

        else:
            self.logger.error(
                "Action malformed (needs 'field' and 'value'): %s" % action
                )
            return False

    def a_remove(self, event, action):
        """
        Remove an event from a field in event or the whole field if no element
        is specified
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

                self.logger.debug("    + %s: Removed: '%s' from '%s'"
                        % (event['rk'], aelement, akey))

            else:
                del event[akey]
                self.logger.debug("    + %s: Removed: '%s'"
                        % (event['rk'], akey))

            return True

        else:
            self.logger.error("Action malformed (needs 'key' and/or 'element'): %s" % action)
            return False


    def a_modify(self, event, action, _name):
        """
        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        derogated = False
        atype = action.get('type', None)
        actionMap = {'override': self.a_override,
                     'remove': self.a_remove}

        if atype in actionMap:
            derogated = actionMap[atype](event, action)

        else:
            self.logger.warning("Unknown action '%s'" % atype)

        # If the event was derogated, fill some informations
        if derogated:
            self.logger.debug("Event changed by rule '%s'" % name)

        return None

    def a_drop(self, event, action, name):
        """
        Drop the event
        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        self.logger.debug("Event dropped by rule '%s'" % name)
        self.drop_event_count += 1
        return DROP

    def a_pass(self, event, action, name):
        """
        Pass the event to the next queue
        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        self.logger.debug("Event passed by rule '%s'" % name)
        self.pass_event_count += 1
        return event

    def a_route(self, event, action, name):
        """
        Change the route to which an event will be sent
        Args:
            event map of the event to be modified
            action map of type action
            _name of the rule
        Returns:
            ``None``
        """

        if "route" in action:
            self.next_amqp_queues = [action["route"]]
            self.logger.debug("Event re-routed by rule '%s'" % name)
        else:
            self.logger.error("Action malformed (needs 'route'): %s" % action)

        return None

    def work(self, event, *xargs, **kwargs):

        rk = get_routingkey(event)
        default_action = self.configuration.get('default_action', 'pass')

        # list of actions supported
        actionMap = {'drop': self.a_drop,
                     'pass': self.a_pass,
                     'override': self.a_modify,
                     'remove': self.a_modify,
                     'route': self.a_route}

        rules = self.configuration.get('rules', [])

        # When list configuration then check black and
        # white lists depending on json configuration
        for filterItem in rules:
            actions = filterItem.get('actions')
            name = filterItem.get('name', 'no_name')

            self.logger.debug(
                'Event: {}, will apply rule {}'.format(event['rk'], filterItem)
                )
            self.logger.debug('filter is {}'.format(filterItem['mfilter']))
            # Try filter rules on current event
            if check(filterItem['mfilter'], event):

                self.logger.debug(
                    'Event: {}, filter matches'.format(event['rk'])
                    )

                for action in actions:
                    if (action['type'] in actionMap):
                        if action['type'] != 'DROP':
                            actions.append({'type': 'pass'})
                        ret = actionMap[action['type']](event, action, name)
                        # If pass then ret == event; end loop
                        if ret:
                            if ret != DROP:
                                event['rk'] = event['_id'] = get_routingkey(event)
                                return event
                            return DROP

                    else:
                        self.logger.warning("Unknown action '%s'" % action)

        # No rules matched
        if default_action == 'drop':
            self.logger.debug("Event '%s' dropped by default action" % (rk))
            self.drop_event_count += 1
            return DROP

        self.logger.debug("Event '%s' passed by default action" % (rk))
        self.pass_event_count += 1

        self.logger.debug('Event before sent to next engine: %s' % event)
        event['rk'] = event['_id'] = get_routingkey(event)

        return event

    def beat(self, *args, **kargs):
        """ Configuration reload for realtime ui changes handling """
        self.derogations = []
        self.configuration = {'rules': [],
                      'default_action': self.find_default_action()}

        self.logger.debug('Reload configuration rules')
        try:
            records = self.storage.find({'crecord_type': 'filter', 'enable': True},
                            sort='priority')

            for record in records:
                record_dump = record.dump()
                self.set_loaded(record_dump)
                record_dump["mfilter"] = literal_eval(record_dump["mfilter"])
                self.logger.debug('Loading record_dump:')
                self.logger.debug(record_dump)
                self.configuration['rules'].append(record_dump)
            self.logger.info('Loaded {} rules'.format(len(self.configuration['rules'])))
            self.send_stat_event()

        except Exception as e:
            self.logger.warning(str(e))


    def set_loaded(self, record):
        if 'run_once' in record and not record['run_once']:
            self.storage.update(record['_id'], {'run_once': True})
            self.logger.info('record {} has been run once'.format(record['_id']))


    def send_stat_event(self):
        """ Send AMQP Event for drop and pass metrics """

        event = forger(
            connector="Engine",
            connector_name="engine",
            event_type="check",
            source_type="resource",
            resource=self.amqp_queue + '_data',
            state=0,
            state_type=1,
            output=("%s event dropped since %s"
                  % (self.drop_event_count, self.beat_interval)),
            perf_data_array = [
                {'metric': 'pass_event',
                 'value': self.pass_event_count,
                 'type': 'GAUGE'},
                {'metric': 'drop_event',
                 'value': self.drop_event_count,
                 'type': 'GAUGE'}])

        self.logger.debug(("%s event dropped since %s"
                   % (self.drop_event_count, self.beat_interval)))
        self.logger.debug(("%s event passed since %s"
                   % (self.pass_event_count, self.beat_interval)))

        rk = get_routingkey(event)
        self.amqp.publish(event, rk, self.amqp.exchange_name_events)
        self.drop_event_count = 0
        self.pass_event_count = 0

    def find_default_action(self):
        """ Find the default action stored and returns it, else assume it default action is pass """

        records = self.storage.find({'crecord_type': 'defaultrule'})
        if records:
            return records[0].dump()["action"]

        self.logger.debug("No default action found. Assuming default action is pass")
        return 'pass'
