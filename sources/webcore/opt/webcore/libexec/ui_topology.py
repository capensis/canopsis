#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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

import sys
from os import listdir
from os.path import expanduser
from logging import getLogger

from bottle import get

logger = getLogger("ui_topology")

operators_path = expanduser('~/opt/amqp2engines/engines/topology')
sys.path.append(operators_path)


@get('/topology/getOperators')
def get_operators():
    operators = []

    for opfile in listdir(operators_path):
        try:
            operator = opfile.split('.')
            if operator[1] == 'py':
                module = __import__(operator[0])
                operators.append(module.options)
                del sys.modules[operator[0]]
        except Exception as err:
            logger.warning("Impossible to parse '%s' (%s)" % (opfile, err))

    return {'total': len(operators), 'success': True, 'data': operators}
