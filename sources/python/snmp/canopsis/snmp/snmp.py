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
from canopsis.old.template import Template
from time import time
from functools import partial
from json import dumps
import re
import socket


manager = RulesManager()
mibs_manager = MibsManager()


class engine(Engine):

    etype = "snmp"

    normal_exchange = 'canopsis.events'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.rules = {}
        # self oids is a cache dict, it only grows with time.
        # It should be cleaned sometime if this engine encounter memory leaks
        self.mibs = {}

    def pre_run(self):
        self.beat()

    def beat(self):

        # load from storage
        self.rules = {rule['oid']['oid']: rule for rule in manager.find()}
        self.logger.info("Loaded {} rules".format(len(self.rules)))
        self.logger.debug(self.rules)

    def work(self, event, *args, **kwargs):

        # this engine works only on snmp trap.
        if event["event_type"] != "trap":
            return

        self.logger.info("Got a trap: {}".format(event))

        # search a rule for the trap OID
        trap_oid = event["snmp_trap_oid"]
        rule = self.rules.get(trap_oid)
        if not rule:
            self.logger.info("No rules for trap {}".format(trap_oid))
            event["snmp_trap_match"] = False
            self.vars_to_string(event)
            self.make_follow(event)
            return

        self.logger.info("Found a rule for trap {}".format(trap_oid))
        errors = []

        # generate a new event
        tt_event = forger(
            connector="Engine",
            connector_name=self.etype,
            event_type="check",
            source_type="resource",
            timestamp=event["timestamp"],
            state_type=0,
        )

        self.logger.debug('start computing template context')
        trap_context = self.get_rule_context(
            rule,
            event.get('snmp_vars', None)
        )

        # if a template exists for any of theses fields, render it.
        for key in ("state", "component", "resource", "output"):
            # is the rule have something we want to change?
            tpl = rule.get(key)
            if tpl is None:
                # rule have no template for this key,
                # just reuse the one in the event if available
                tt_event[key] = event.get(key)
                continue

            # a template has been found for the key
            # do the rendering!
            # meaning, replace the {{oid}} with the vars of the snmp trap
            self.error = None

            try:
                # Try to convert the template to unicode
                unicode_template = unicode(str(tpl).decode('utf-8'))
                value = Template(unicode_template)(trap_context)
            except Exception as e:
                self.logger.error(
                    'Error, encoding problem in this field {}'.format(e)
                )
                value = tpl

            self.logger.debug(
                '"{}" field had template "{}" set to "{}"'.format(
                    key,
                    tpl,
                    value
                )
            )

            tt_event[key] = value
            if self.error is not None:
                errors.append("{}: {}".format(key, self.error))

        # and show that the event got a match in our trap/rules db.
        event["snmp_trap_match"] = True
        if errors:
            event["snmp_trap_errors"] = errors
            self.logger.info("No new events, trap as errors: {}".format(
                errors
            ))
        else:
            # publish the new event
            self.logger.info("Publish a new event: {}".format(tt_event))
            rk = get_routingkey(tt_event)
            self.amqp.publish(
                tt_event, rk, self.amqp.exchange_name_events)

        # Allow mongo json with dotted key values persistance
        self.vars_to_string(event)

        self.make_follow(event)

    def get_rule_context(self, rule, snmp_vars):
        # Computes the template context from event information
        # and snmp rules information
        if snmp_vars is None:
            self.logger.debug('no snmp vars in event')
            return {}
        else:
            context = {}
            mib = self.get_and_cache_mib(rule)
            # When unable to retrieve mib information
            if mib is None:
                self.logger.debug('no mib info found')
                return context
            else:
                for mibobject in mib['objects']:
                    context[mibobject] = snmp_vars.get(mib['oid'], '')
                self.logger.debug('generated context {}'.format(context))
                return context

    def get_and_cache_mib(self, rule):

        _id = '{}::{}'.format(
            rule['oid']['moduleName'],
            rule['oid']['mibName']
        )

        if _id in self.mibs:
            self.logger.debug('mib found in cache')
            return self.mibs[_id]
        else:
            result = list(mibs_manager.get(oids=[_id]))
            oid = None
            if len(result):
                oid = result[0]['oid']

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

            #Test oid TODO remove
            oid = '1.3.6.1.4.1.20006.1.3.1.17'
            if oid is not None and objects is not None:

                self.mibs[_id] = {
                    'oid': oid,
                    'objects': objects
                }
                return self.mibs[_id]
            else:
                return None

    def vars_to_string(self, event):
        key = 'snmp_vars'
        if key in event and isinstance(event[key], dict):
            event[key] = dumps(event[key])

    def make_follow(self, event):
        rk = get_routingkey(event)
        event['_id'] = rk
        self.amqp.publish(
            event, rk, self.normal_exchange)
