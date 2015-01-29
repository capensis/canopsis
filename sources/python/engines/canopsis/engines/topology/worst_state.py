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
name = "worst_state"
display_name = "Worst State"
description = "Calcul the worst state"

options = {
    '_id': name,
    'component': display_name,
    'description': description,
    'event_type': 'operator',
    'source_type': 'component',
    'nodeMaxOutConnexion': 10,
    'nodeMaxInConnexion': 10
}


def operator(states, options={}):
    logger.debug("%s: Calcul state for %s" % (name, states))

    if len(states) == 1:
        state = states[0]
    else:
        states.sort()
        states.reverse()
        state = states[0]

    logger.debug("%s: + State: %s" % (name, state))
    return state
