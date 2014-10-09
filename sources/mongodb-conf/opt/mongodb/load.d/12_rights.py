#!/usr/bin/env python
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

from canopsis.old.account import Account
from canopsis.organisation.rights import Rights

from os.path import join
import json
import sys

root = Account(user="root", group="root")
right_module = Rights()
actions_path = join(sys.prefix, 'opt/mongodb/load.d/rights/actions_ids.json')

def init():
    json_data = open(actions_path)
    data = json.load(json_data)

    for action_id in data:
        right_module.add(action_id,
                         data[action_id].get('desc', "Empty desc"))

def update():
    pass
