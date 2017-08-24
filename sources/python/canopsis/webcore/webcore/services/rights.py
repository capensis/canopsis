# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 'Capensis' [http://www.capensis.com]
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

from canopsis.organisation.rights import Rights
from canopsis.common.ws import route
from bottle import HTTPError


rights = None


def get_manager():
    global rights

    if not rights:
        rights = Rights()

    return rights


def save_group(ws, group):
    cname = group['_id']
    crights = group['rights']

    group = rights.get_group(cname)

    if not group and not rights.create_group(cname, crights):
        raise ws.Error('Impossible to create group')

    rights.update_rights(cname, 'group', crights, group)

    return group


def save_profile(ws, profile):
    pid = profile['_id']
    pgroup = profile['profile_groups']
    prights = profile['profile_rights']

    profile = rights.get_profile(pid)

    if not profile and not rights.create_profile(pid, pgroup):
        raise ws.Error('Impossible to create profile')

    rights.update_group(pid, 'profile', pgroup, profile)
    rights.update_rights(pid, 'profile', prights, profile)

    return profile


def save_role(ws, role):
    role_payload = role
    rid = role['_id']

    rgroup = role.get('groups', None)
    rrights = role.get('rights', None)
    rprofile = role.get('profile', None)

    defaultview = role.get('defaultview', None)

    role = rights.get_role(rid)

    if defaultview is not None:
        rights.update_fields(rid, 'role', {'defaultview': defaultview})

    if not role and not rights.create_role(rid, rprofile):
        raise ws.Error('Impossible to create role')

    if rprofile:
        rights.update_profile(rid, 'role', rgroup, role)

    if rgroup:
        rights.update_group(rid, 'role', rgroup, role)

    if rrights:
        rights.update_rights(rid, 'role', rrights, role)

    return role_payload


def save_user(ws, record):
    uid = record.pop('_id')
    urole = record.pop('role')
    ucontact = record.pop('contact', None)
    urights = record.pop('rights', None)
    ugroup = record.pop('groups', None)

    # Remove authkey from record when None
    # Prevents update with None value
    record.pop('authkey', None)

    if ucontact is None:
        ucontact = {
            'name': '{0} {1}'.format(
                record.get('firstname', ''),
                record.get('lastname', ''),
            ),
            'email': record.get('mail', '')
        }

    user = rights.get_user(uid)

    if not user:
        user = rights.create_user(
            uid, urole,
            contact=ucontact,
            rights=urights,
            groups=ugroup
        )

        if not user:
            raise ws.Error('Impossible to create user')

    if ugroup is not None:
        rights.update_group(uid, 'user', ugroup, user)

    if urights is not None:
        rights.update_rights(uid, 'user', urights, user)

    rights.update_fields(uid, 'user', record)

    if not rights.add_role(uid, urole):
        raise ws.Error('Impossible to add user to role')

    return record


def exports(ws):
    mgr = get_manager()

    @route(ws.application.get)
    def rights(uid):
        urights = mgr.get_user_rights(uid)

        if not urights:
            raise HTTPError(404, 'No rights found for user: {0}'.format(uid))

        else:
            return urights

    @route(ws.application.post, name='account/group', payload=['group'])
    def create_group(group):
        return save_group(ws, group)

    @route(ws.application.put, name='account/group', payload=['group'])
    def update_group(_id, group):
        return save_group(ws, group)

    @route(ws.application.post, name='account/profile', payload=['profile'])
    def create_profile(profile):
        return save_profile(ws, profile)

    @route(ws.application.put, name='account/profile', payload=['profile'])
    def update_profile(_id, profile):
        return save_profile(ws, profile)

    @route(ws.application.post, name='account/role', payload=['role'])
    def create_role(role):
        return save_role(ws, role)

    @route(ws.application.put, name='account/role', payload=['role'])
    def update_role(_id, role):
        return save_role(ws, role)

    @route(ws.application.post, name='account/user', payload=['user'])
    def create_user(user):
        return save_user(ws, user)

    @route(ws.application.put, name='account/user', payload=['user'])
    def update_user(_id, user):
        return save_user(ws, user)

    @route(ws.application.delete, name='account/delete')
    def delete_entity(etype, _id):
        if not mgr.delete(etype, _id):
            raise ws.Error('Unknown entity type: {0}'.format(etype))
