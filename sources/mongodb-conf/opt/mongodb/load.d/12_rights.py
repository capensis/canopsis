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
role_path = join(sys.prefix, 'opt/mongodb/load.d/rights/default_roles.json')
group_path = join(sys.prefix, 'opt/mongodb/load.d/rights/default_groups.json')
profile_path = join(sys.prefix, 'opt/mongodb/load.d/rights/default_profiles.json')

def add_actions(data):
    for action_id in data:
        right_module.add(action_id,
                         data[action_id].get('desc', "Empty desc"))

def add_users(data):
    for user in data:
        right_module.delete_user(user['_id'])
        right_module.create_user(user['_id'], user.get('role', None),
                                 rights=user.get('rights', None),
                                 contact=user.get('contact', None),
                                 groups=user.get('groups', None))

def add_roles(data):
    for role in data:
        right_module.create_role(role['_id'], role.get('profile', None))
        record = right_module.get_role(role['_id'])
        right_module.update_rights(
            role['_id'], 'role', role.get('rights', {}), record
            )
        right_module.update_group(
            role['_id'], 'role', role.get('groups', []), record
            )

def add_profiles(data):
    for profile in data:
        right_module.create_profile(profile['_id'], None)
        record = right_module.get_profile(profile['_id'])
        right_module.update_rights(
            profile['_id'], 'profile', profile.get('rights', {}), record
            )
        right_module.update_group(
            profile['_id'], 'profile', profile.get('groups', []), record
            )

def add_groups(data):
    for group in data:
        right_module.create_group(group['_id'], None)
        record = right_module.get_group(group['_id'])
        right_module.update_rights(
            group['_id'], 'group', group.get('rights', {}), record
            )

def init():
    add_actions(json.load(open(actions_path)))
    add_users(json.load(open(user_path)))
    add_roles(json.load(open(role_path)))
    # add_groups(json.load(open(group_path)))
    # add_profiles(json.load(open(profile_path)))

def update():
    init()
