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
from time import time
from functools import partial
from json import dumps
import re
import socket

manager = RulesManager()
mibs_manager = MibsManager()

TEMPLATE_PATTERN = re.compile(r"\{\{([^}]+)\}\}")


class engine(Engine):

    etype = "snmp"

    normal_exchange = 'canopsis.events'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        self.funcmap = {
            "resolveip": self._func_resolveip,
            "upper": self._func_upper,
            "lower": self._func_lower
        }
        self.rules = {}
        # self oids is a cache dict, it only grows with time.
        # It should be cleaned sometime if this engine encounter memory leaks
        self.oids = {}

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

        # prepare the templating function
        self.logger.info("Found a rule for trap {}".format(trap_oid))
        errors = []
        f_repl = partial(self._template_repl, rule, event["snmp_vars"])

        # generate a new event
        tt_event = forger(
            connector="Engine",
            connector_name=self.etype,
            event_type="check",
            source_type="resource",
            timestamp=event["timestamp"],
            state_type=0,
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
            value = re.sub(TEMPLATE_PATTERN, f_repl, tpl)
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

    def get_and_cache_oid(self, _id):

        if _id in self.oids:
            return self.oids[_id]
        else:
            result = list(mibs_manager.get(ids=[_id]))
            if len(result):
                oid = result['oid']
            else:
                oid = None
            self.oids[_id] = oid
            return oid

    def vars_to_string(self, event):
        key = 'snmp_vars'
        if key in event and isinstance(event[key], dict):
            event[key] = dumps(event[key])

    def make_follow(self, event):
        rk = get_routingkey(event)
        event['_id'] = rk
        self.amqp.publish(
            event, rk, self.normal_exchange)

    def _template_repl(self, rule, vars, matchobj):

        m = matchobj.group(1)
        items = [item.strip() for item in m.split("|")]

        # Find oids the way data are stored
        _id = '{}::{}'.format(rule['moduleName'], rule['mibName'])
        oid = self.get_and_cache_oid(_id)

        value = vars.get(oid)

        if value is None:
            self.error = "variable {} missing".format(oid)
            return ""

        for funcname in items[1:]:
            f = self.funcmap.get(funcname)
            if not f:
                self.error = "unknown function {}".format(funcname)
                return
            value = f(value)

        return value

    def _func_resolveip(self, value):
        # dns resolve
        try:
            hostname = socket.gethostbyaddr(value)
        except:
            pass
        else:
            value = hostname[0]
        return value

    def _func_upper(self, value):
        return value.upper()

    def _func_lower(self, value):
        return value.lower()
