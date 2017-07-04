# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals
from six import string_types
from canopsis.common.ws import route
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.watcher.manager import Watcher

watcher_manager = Watcher()

def check_values(data):
    """Check if the values present in data respect the specification. If
    the values are correct do nothing. If not, raise a ValueError
    :param dict data: the data."""

    def check(data, key, type_):
        if key in data and data[key] is not None:
            if not isinstance(data[key], type_):
                raise ValueError("The {0} must be a {1} not {2}".format(
                    key, type_, type(data[key])))

    # check str values
    for k in ["name", "author", "rrule", "component", "connector",
              "connector_name"]:
        check(data, k, string_types)

    # check int values
    for k in ["tstart", "tstop"]:
        check(data, k, int)

    # check dict values
    for k in ["comments"]:

        if "comments" not in data:
            continue

        if data["comments"] is None:
            continue

        check(data, k, list)

        for elt in data["comments"]:
            if not isinstance(elt, dict):
                raise ValueError("The list {0} store only {1} not {2}".format(
                    k, dict, type(elt)))

    if "enabled" not in data or data["enabled"] is None or isinstance(data['enabled'], bool):
        return

    if data["enabled"] in ["True", "true"]:
        data["enabled"] = True
    elif data["enabled"] in ["False", "false"]:
        data["enabled"] = False
    else:
        raise ValueError("The enabled value does not match a boolean")


def exports(ws):

    pbm = PBehaviorManager()

    @route(
        ws.application.post,
        name='pbehavior/create',
        payload=[
            'name', 'filter', 'author',
            'tstart', 'tstop', 'rrule',
            'enabled', 'comments',
            'connector', 'connector_name'
        ]
    )
    def create(
            name, filter, author,
            tstart, tstop, rrule=None,
            enabled=True, comments=None,
            connector='canopsis', connector_name='canopsis'
    ):

        data = locals()
        check_values(data)

        result =  pbm.create(
            name=name, filter=filter, author=author,
            tstart=tstart, tstop=tstop, rrule=rrule,
            enabled=enabled, comments=comments,
            connector=connector, connector_name=connector_name
        )

        watcher_manager.compute_watchers()

        return result

    @route(
        ws.application.get,
        name='pbehavior/read',
        payload=['_id']
    )
    def read(_id=None):
        ok = False
        if isinstance(_id, string_types):
            ok = True
        elif isinstance(_id, list):
            ok = True
            for elt in _id:
                if not isinstance(elt, string_types):
                    ok = False

        if not ok:
            raise ValueError("_id should be str, a list, None (null) not"\
                             "{0}".format(type(_id)))

        return pbm.read(_id)

    @route(
        ws.application.put,
        name='pbehavior/update',
        payload=[
            '_id',
            'name', 'filter',
            'tstart', 'tstop', 'rrule',
            'enabled'
        ]
    )
    def update(
            _id,
            name=None, filter=None,
            tstart=None, tstop=None, rrule=None,
            enabled=None, comments=None,
            connector=None, connector_name=None,
            author=None
    ):
        params = locals()
        check_values(params)
        params.pop('_id')
        params.pop('pbm')

        return pbm.update(_id, **params)

    @route(
        ws.application.delete,
        name='pbehavior/delete',
        payload=['_id']
    )
    def delete(_id):

        return pbm.delete(_id)

    @route(
        ws.application.post,
        name='pbehavior/comment/create',
        payload=['pbehavior_id', 'author', 'message']
    )
    def create_comment(pbehavior_id, author, message):
        author = str(author)
        message = str(author)
        return pbm.create_pbehavior_comment(pbehavior_id, author, message)

    @route(
        ws.application.put,
        name='pbehavior/comment/update',
        payload=['pbehavior_id', '_id', 'author', 'message']
    )
    def update_comment(pbehavior_id, _id, author=None, message=None):
        author = str(author)
        message = str(message)
        return pbm.update_pbehavior_comment(
            pbehavior_id, _id,
            author=author, message=message
        )

    @route(
        ws.application.delete,
        name='pbehavior/comment/delete',
        payload=['pbehavior_id', '_id']
    )
    def delete_comment(pbehavior_id, _id):
        return pbm.delete_pbehavior_comment(pbehavior_id, _id)
