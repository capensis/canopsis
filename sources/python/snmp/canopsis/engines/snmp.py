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
from canopsis.check.manager import CheckManager
from canopsis.context.manager import Context
from canopsis.event import get_routingkey, forger
from time import time
from json import loads
from functools import partial
import re

context_manager = Context()
check_manager = CheckManager()

TEMPLATE_PATTERN = re.compile(r"\{\{([\d\.]+)\}\}")

RULES = {
    "1.3.6.1.6.3.1.1.5.3": {
        "component": None,
        "resource": "link:{{1.3.6.1.2.1.2.2.1.1}}",
        "state": "1",
        "output": "Link down on interface {{1.3.6.1.2.1.2.2.1.1}}"
    },
    "1.3.6.1.6.3.1.1.5.4": {
        "component": None,
        "resource": "link:{{1.3.6.1.2.1.2.2.1.1}}",
        "state": "0",
        "output": "Link up on interface {{1.3.6.1.2.1.2.2.1.1}}"
    }
}


class engine(Engine):
    etype = "snmp"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        self.rules = {}
        self.beat()

    def beat(self):
        # TODO: load rules from database
        self.rules = RULES

    def work(self, event, *args, **kwargs):
        # this engine works only on snmp trap.
        if event["event_type"] != "trap":
            return

        # load the snmp serialized message embedded in the output
        # (done by snmp2canopsis connector)
        snmp = loads(event["output"])
        self.logger.info("Got a trap: {}".format(snmp))

        # search a rule for the trap OID
        trap_oid = snmp["trap_oid"]
        rule = self.rules.get(trap_oid)
        if not rule:
            self.logger.info("No rules for trap {}".format(trap_oid))
            event["snmp_trap_match"] = False
            return event

        # prepare the templating function
        self.logger.info("Found a rule for trap {}".format(trap_oid))
        f_repl = partial(self._template_repl, rule, snmp.get("vars"))

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
            value = re.sub(TEMPLATE_PATTERN, f_repl, tpl)
            tt_event[key] = value

        # publish the new event
        self.logger.info("Publish a new event: {}".format(tt_event))
        rk = get_routingkey(tt_event)
        self.amqp.publish(
            tt_event, rk, self.amqp.exchange_name_events)

        # and show that the event got a match in our trap/rules db.
        event["snmp_trap_match"] = True
        return event

    def _template_repl(self, rule, vars, matchobj):
        oid = matchobj.group(1)
        return vars.get(oid)
