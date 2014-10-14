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
user_path = join(sys.prefix, 'opt/mongodb/load.d/rights/default_users.json')

def add_actions(data):
    for action_id in data:
        right_module.add(action_id,
                         data[action_id].get('desc', "Empty desc"))

def add_users(data):
    for user in data:
        right_module.delete_user(user['_id'])
        right_module.create_user(user['_id'], user.setdefault('role', None),
                                 rights=user.setdefault('rights', None),
                                 contact=user.setdefault('contact', None),
                                 groups=user.setdefault('groups', None))

def init():
    add_actions(json.load(open(actions_path)))
    add_users(json.load(open(user_path)))

def update():
    init()
