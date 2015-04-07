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
from canopsis.snmp.manager import SnmpManager
from canopsis.context.manager import Context
from canopsis.event import get_routingkey, forger
from time import time
from functools import partial
import re
import socket

manager = SnmpManager()

TEMPLATE_PATTERN = re.compile(r"\{\{([^}]+)\}\}")

RULES = {
    "1.3.6.1.6.3.1.1.5.3": {
        "component": None,
        "resource": "link:{{1.3.6.1.2.1.2.2.1.1}}",
        "state": "1",
        "output": "Link down on interface {{1.3.6.1.2.1.2.2.1.1}}"
    },
    "1.3.6.1.6.3.1.1.5.4": {
        "component": "{{1.3.6.1.2.1.1.5|resolveip}}",
        "resource": "link:{{1.3.6.1.2.1.2.2.1.1}}",
        "state": "0",
        "output": "Link up on interface {{1.3.6.1.2.1.2.2.1.1}}"
    }
}


class engine(Engine):
    etype = "snmp"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        self.funcmap = {
            "resolveip": self._func_resolveip,
            "upper": self._func_upper,
            "lower": self._func_lower
        }
        self.rules = {}
        self.beat()

    def beat(self):
        # load in storage
        self.rules = {rule["_id"]: rule for rule in manager.get()}
        self.logger.info("Loaded {} rules".format(len(self.rules)))

        if not self.rules:
            # FIXME: put some data into the db before using it.
            self.logger.info("Insert default rules to starts with.")
            for oid, rule in RULES.items():
                manager.put(oid, rule)
            self.rules = RULES
            self.logger.info("Inserted {} rules".format(len(self.rules)))

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
            return event

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
            self.logger.info("No new events, trap as errors: {}".format(errors))
        else:
            # publish the new event
            self.logger.info("Publish a new event: {}".format(tt_event))
            rk = get_routingkey(tt_event)
            self.amqp.publish(
                tt_event, rk, self.amqp.exchange_name_events)

        return event

    def _template_repl(self, rule, vars, matchobj):
        m = matchobj.group(1)
        items = [item.strip() for item in m.split("|")]
        oid = items[0]
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
