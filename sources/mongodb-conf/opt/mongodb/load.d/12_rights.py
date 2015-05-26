#!/usr/bin/env python
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

from canopsis.old.account import Account
from canopsis.organisation.rights import Rights

from os.path import join
import json
import sys
import uuid

root = Account(user="root", group="root")
right_module = Rights()


def load(path):
    absolute_path = join(
        sys.prefix,
        path
    )
    try:
        return json.load(open(absolute_path))
    except Exception as e:
        print('\n\t/!\ Malformed right management json file : {}\n'.format(
            absolute_path
        ))
        sys.exit(1)

actions = load('opt/mongodb/load.d/rights/actions_ids.json')
users = load('opt/mongodb/load.d/rights/default_users.json')
roles = load('opt/mongodb/load.d/rights/default_roles.json')


def add_actions(data, clear):
    for action_id in data:
        if right_module.get_action(action_id) is None or clear:

            print('Action initialization : {}'.format(action_id))

            right_module.add(action_id,
                             data[action_id].get('desc', "Empty desc"))


def add_users(data, clear):
    for user in data:
        if right_module.get_user(user['_id']) is None or clear:

            print('User initialization : {}'.format(user['_id']))

            right_module.create_user(user['_id'], user.get('role', None),
                                     rights=user.get('rights', None),
                                     contact=user.get('contact', None),
                                     groups=user.get('groups', None))

            right_module.update_fields(user['_id'], 'user', {
                'external': user.get('external', False),
                'enable': user.get('enable', True),
                'shadowpasswd': user.get('shadowpass', None),
                'authkey': str(uuid.uuid1())
            })


def add_roles(data, clear):
    for role in data:

        if right_module.get_role(role['_id']) is None or clear:

            print('Role initialization : {}'.format(role['_id']))

            right_module.create_role(role['_id'], role.get('profile', None))
            record = right_module.get_role(role['_id'])
            right_module.update_rights(
                role['_id'], 'role', role.get('rights', {}), record
                )
            right_module.update_group(
                role['_id'], 'role', role.get('groups', []), record
                )


def init(clear=True):
    add_actions(actions, True)
    add_users(users, clear)
    add_roles(roles, clear)


def update():
    init(clear=False)
