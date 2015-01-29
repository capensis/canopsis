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

# Fill by engine
logger = None
name = "and"
display_name = "And"
description = "Apply AND between states"

options = {
    '_id': name,
    'component': display_name,
    'description': description,
    'event_type': 'operator',
    'source_type': 'component',
    'nodeMaxOutConnexion': 10,
    'nodeMaxInConnexion': 10,
    'form': {
        'xtype': 'form',
        'items': [
            {
                "xtype": "combobox",
                "name": "state",
                "fieldLabel": "Check state",
                "queryMode": "local",
                "displayField": "text",
                "valueField": "value",
                "value": "0",
                "store": {
                    "xtype": "store",
                    "fields": ["value", "text"],
                    "data": [
                        {"value": "-1", "text": "Not Ok"},
                        {"value": "0", "text": "Ok"},
                        {"value": "1", "text": "Warning"},
                        {"value": "2", "text": "Critical"},
                        {"value": "3", "text": "Unknown"}
                    ]
                }
            }, {
                "xtype": "combobox",
                "name": "Then",
                "fieldLabel": "then",
                "queryMode": "local",
                "displayField": "text",
                "valueField": "value",
                "value": "0",
                "store": {
                    "xtype": "store",
                    "fields": ["value", "text"],
                    "data": [
                        {"value": "-1", "text": "Worst State"},
                        {"value": "0", "text": "Ok"},
                        {"value": "1", "text": "Warning"},
                        {"value": "2", "text": "Critical"},
                        {"value": "3", "text": "Unknown"}
                    ]
                }
            }, {
                "xtype": "combobox",
                "name": "else",
                "fieldLabel": "Else",
                "queryMode": "local",
                "displayField": "text",
                "valueField": "value",
                "value": "-1",
                "store": {
                    "xtype": "store",
                    "fields": ["value", "text"],
                    "data": [
                        {"value": "-1", "text": "Worst State"},
                        {"value": "0", "text": "Ok"},
                        {"value": "1", "text": "Warning"},
                        {"value": "2", "text": "Critical"},
                        {"value": "3", "text": "Unknown"}
                    ]
                }
            }
        ]
    }
}


def operator(states, options={}):
    logger.debug("%s: Calcul state for %s" % (name, states))

    and_state = int(options.get('state', 0))
    and_then = int(options.get('then', 0))
    and_else = int(options.get('else', -1))

    state = 3

    if and_then == -1 or and_else == -1:
        states.sort()
        states.reverse()
        worst_state = states[0]

    if len(states) == len([s for s in states if s == and_state]):
        if and_then == -1:
            state = worst_state
        else:
            state = and_then
    else:
        if and_else == -1:
            state = worst_state
        else:
            state = and_else

    logger.debug("%s: + State: %s" % (name, state))
    return state
