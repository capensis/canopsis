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

from canopsis.engines.core import Engine
from canopsis.snmp.rulesmanager import RulesManager
from canopsis.snmp.mibs import MibsManager
from canopsis.context.manager import Context
from canopsis.event import get_routingkey, forger
from canopsis.common.template import Template
from time import time
from functools import partial
from json import dumps
import re
import socket
from time import time

manager = RulesManager()
mibs_manager = MibsManager()


class engine(Engine):

    etype = 'snmp'

    normal_exchange = 'canopsis.events'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.rules = {}

    def pre_run(self):
        self.beat()

    def beat(self):

        """
        On beat the engine gets all snmp rules from database
        and put them in memory as a dict with key equals an oid
        """
        # Load snmp rules from database
        self.rules = {}
        for rule in manager.find():
            oid = rule.get('oid', None).get('oid', None)
            if oid is not None:
                self.rules[oid] = rule

        self.logger.info("Loaded {} rules".format(len(self.rules)))
        self.logger.debug(self.rules)

    def work(self, event, *args, **kwargs):

        """
        When a trap event is received, if a rule matches it s oid,
        then the rule describes event component, resource, output and state
        transformation. If the transformation succeeds the translated event is
        sent to the normal canopsis event exchange, otherwise, the trap is
        tagged as untranslated and is sent to the event log collection
        """

        # This engine works only on snmp trap.
        if event['event_type'] != 'trap':
            return

        self.logger.debug('Got a trap event: {}'.format(event))

        # Search a rule for the trap OID
        trap_oid = event['snmp_trap_oid']

        rule = self.rules.get(trap_oid)
        if not rule:
            message = 'No rules for trap {}'.format(trap_oid)
            return self.on_trap_error(event, message)

        self.logger.info('Found a rule for trap {}'.format(trap_oid))
        errors = []

        # generate a new event
        translated_event = forger(
            connector='Engine',
            connector_name=self.etype,
            event_type='check',
            source_type='resource',
            timestamp=event['timestamp'],
            state_type=0,
        )

        self.logger.debug('Start computing template context')
        trap_context = self.get_trap_context(
            rule,
            event.get('snmp_vars', None),
            errors
        )

        if errors:
            message = 'Unable to get context for trap {}'.format(trap_oid)
            return self.on_trap_error(event, message, errors)

        # if a template exists for any of theses fields, render it.
        for key in ('state', 'component', 'resource', 'output'):
            # is the rule have something we want to change?
            template = rule.get(key)
            if template is None:
                # rule have no template for this key,
                # just reuse the one in the event if available
                translated_event[key] = event.get(key)
                message = 'Key not managed in rule {}'.format(key)
                return self.on_trap_error(event, message)

            # a template has been found for the key
            # do the rendering!
            # meaning, replace the {{oid}} with the vars of the snmp trap

            try:
                # Try to convert the template to unicode
                unicode_template = unicode(str(template).decode('utf-8'))
                value = Template(unicode_template)(trap_context)
            except Exception as e:
                message = 'Key {}, Template {}: {}'.format(key, template, e)
                return self.on_trap_error(event, message)

            if not value.strip():
                message = 'Empty key value : {}'.format(key)
                return self.on_trap_error(event, message)

            if key == 'state' and value.strip() not in ['0', '1', '2', '3']:
                message = 'Invalid state : {}'.format(value)
                return self.on_trap_error(event, message)
            elif key == 'state':
                # Parse state to int
                value = int(value)

            self.logger.debug(
                '"{}" field had template "{}" set to "{}"'.format(
                    key,
                    template,
                    value
                )
            )

            translated_event[key] = value

        return self.on_trap_translated(translated_event)

    def on_trap_translated(self, event):

        """
        :param: event, the translated event to send to canopsis exchange
        When trap translation succeeds
        """

        # Show that the event got a match in our trap/rules db.
        event["snmp_trap_match"] = True
        self.make_follow(event)
        return event

    def on_trap_error(self, event, reason, errors=None):

        """
        :param: event the untranslated event to send to canopsis exchange
        :param: reason is a string that describes why translation is wrong
        :param: errors can be a list of string that is put in the sent event
        When trap translation fails, event is tagged and
        error information is set to the event.
        """

        if errors is None:
            errors = []

        self.logger.info(reason)
        event["snmp_trap_match"] = False
        errors.append(reason)
        event['errors'] = errors
        self.make_follow(event)
        return event

    def get_trap_context(self, rule, snmp_vars, errors):

        """
        :param: rule is the current event matching rule as dict
        :param: snmp_vars are oid information holded by current event
        :param: errors is a list of string error to fill when error occurs
        Fetch data base from specific rule and event information in order to
        replace oid values by objects names from mib definition
        for human readability
        """

        # Computes the template context from event information
        # and snmp rules information
        if snmp_vars is None:
            message = 'No snmp vars in event'
            errors.append(message)
            self.logger.debug(message)
            return None
        else:
            objects = self.get_mib(rule)
            # When unable to retrieve mib information
            if objects is None:
                message = 'No mib info found'
                errors.append(message)
                self.logger.error(message)
                return None
            else:
                context = self.get_mibs_objects(
                    rule,
                    objects,
                    snmp_vars,
                    errors
                )
                self.logger.debug('generated context {}'.format(context))
                return context

    def get_mibs_objects(self, rule, mib_objects, snmp_vars, errors):

        """
        :param: rule is the current event matching rule as dict
        :param: mib_objects are objects name in mib database for a given module
        :param: snmp_vars are oid information holded by current event
        :param: errors is a list of string error to fill when error occurs
        Same as get_trap_context method, for mib object information
        """

        # Data validation
        if mib_objects is None:
            message = 'Mib does not contains objects'
            errors.append(message)
            self.logger.error(message)
            return None

        # Template context generation and mib module objects caching
        context = {}
        for mib_object in mib_objects:

            _id = '{}::{}'.format(
                rule['oid']['moduleName'],
                mib_object
            )

            result = list(mibs_manager.get(oids=[_id]))

            if len(result):
                oid = result[0]['oid']
                context[mib_object] = snmp_vars.get(oid, '')
            else:
                errors.append('Mib object oid not found in db : {}'.format(
                    _id
                ))
                context[mib_object] = None

        return context

    def get_mib(self, rule):

        """
        :param: rule is the current event matching rule as dict
        Fetch mib information from storage from a matchin rule
        """

        _id = '{}::{}'.format(
            rule['oid']['moduleName'],
            rule['oid']['mibName']
        )

        result = list(mibs_manager.get(oids=[_id]))
        oid = None
        if len(result):
            oid = result[0]['oid']

        # Query storage from mib identifiers
        query = {
            'moduleName': rule['oid']['moduleName'],
            'name': rule['oid']['mibName']
        }

        result = list(mibs_manager.get(query=query))
        objects = None
        if len(result):
            objects = result[0]['objects'].keys()

        self.logger.debug('fetch oid {}, objects {}'.format(
            oid,
            len(objects)
        ))

        return objects

    def make_follow(self, event):

        """
        Sends an event to the normal canopsis exchange as event
        from this engine are bind to an alternative exchange
        """

        # Allow mongo json with dotted key values persistance
        key = 'snmp_vars'
        if key in event and isinstance(event[key], dict):
            event[key] = dumps(event[key])

        # Publish event
        rk = get_routingkey(event)
        event['_id'] = rk
        self.amqp.publish(
            event, rk, self.normal_exchange
        )
