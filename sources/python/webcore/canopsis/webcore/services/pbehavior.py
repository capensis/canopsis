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
from json import loads
from six import string_types

from canopsis.common.ws import route
from canopsis.pbehavior.utils import check_valid_rrule
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.watcher.manager import Watcher as WatcherManager

def check(data, key, type_):
    """
    If key exists in data and data[key] != None, check for type_
    :param dict data: data to check
    :param str key: key in data dict
    :param type type_: data[key] type to check
    :raises ValueError: type_ doesn't match value's type
    """
    if key in data and data[key] is not None:
        if not isinstance(data[key], type_):
            raise ValueError("The {0} must be a {1} not {2}".format(
                key, type_, type(data[key])))

def check_values(data):
    """Check if the values present in data respect the specification. If
    the values are correct do nothing. If not, raises an error.
    :raises ValueError: a value is invalid.
    :param dict data: the data."""

    # check str values
    for k in ["name", "author", "rrule", "component", "connector",
              "connector_name"]:
        check(data, k, string_types)

    # check int values
    for k in ["tstart", "tstop"]:
        check(data, k, int)

        if 'tstart' in data and 'tstop' in data:
            if data['tstart'] >= data['tstop']:
                raise ValueError('tstop cannot be inferior or equal to tstart')

    # check dict values
    for k in ["comments"]:

        if "comments" not in data:
            continue

        if data["comments"] is None:
            continue

        check(data, k, list)

        for elt in data["comments"]:
            if not isinstance(elt, dict):
                raise ValueError("The list {0} store only {1} not {2}"
                                 .format(k, dict, type(elt)))

    if "filter" in data and isinstance(data["filter"], string_types):
        try:
            data["filter"] = loads(data["filter"])
        except ValueError:
            raise ValueError("Cant decode mfilter parameter: {}"
                             .format(data["filter"]))

    if 'rrule' in data:
        check_valid_rrule(data['rrule'])

    # useful when enabled doesn't exist in document
    if "enabled" not in data or data["enabled"] is None or isinstance(data['enabled'], bool):
        return

    if data["enabled"] in ["True", "true"]:
        data["enabled"] = True
    elif data["enabled"] in ["False", "false"]:
        data["enabled"] = False
    else:
        raise ValueError("The enabled value does not match a boolean")

class RouteHandlerPBehavior(object):

    def __init__(self, pb_manager, watcher_manager):
        """
        :param PBehaviorManager pb_manager: pbehavior manager
        :param WatcherManager watcher_manager: watcher manager
        """
        super(RouteHandlerPBehavior, self).__init__()
        self.watcher_manager = watcher_manager
        self.pb_manager = pb_manager

    def create(self, name, filter_, author,
               tstart, tstop, rrule,
               enabled, comments,
               connector, connector_name):
        data = locals()
        data.pop('self')
        check_values(data)

        result = self.pb_manager.create(
            name=name, filter=filter_, author=author,
            tstart=tstart, tstop=tstop, rrule=rrule,
            enabled=enabled, comments=comments,
            connector=connector, connector_name=connector_name
        )

        self.watcher_manager.compute_watchers()

        return result

    def read(self, _id):
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

        return self.pb_manager.read(_id)

    def update(self, _id, name, filter, tstart, tstop, rrule,
            enabled, comments, connector, connector_name, author):
        params = locals()
        params.pop('_id')
        params.pop('self')
        check_values(params)

        return self.pb_manager.update(_id, **params)

    def delete(self, _id):
        return self.pb_manager.delete(_id)

    def create_comment(self, pb_id, author, message):
        author = str(author)
        message = str(author)
        return self.pb_manager.create_pbehavior_comment(pbehavior_id, author, message)

    def update_comment(self, pb_id, _id, author, message):
        author = str(author)
        message = str(message)
        return self.pb_manager.update_pbehavior_comment(
            pb_id, _id, author=author, message=message)

    def delete_comment(self, pbehavior_id, _id):
        return self.pb_manager.delete_pbehavior_comment(pbehavior_id, _id)

def exports(ws):

    rhpb = RouteHandlerPBehavior(
        pbm=PBehaviorManager(), watcher_manager=WatcherManager())

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

        return rhpb.create(
            name, filter, author, tstart, tstop,
            rrule, enabled, comments, connector, connector_name)

    @route(
        ws.application.get,
        name='pbehavior/read',
        payload=['_id']
    )
    def read(_id=None):
        return rhpb.read(_id)

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
        return rhpb.update(
            _id, name=name, filter_=filter, tstart=tstart, tstop=tstop,
            rrule=rrule, enabled=enabled, comments=comments,
            connector=connector, connector_name=connector_name, author=author)


    @route(
        ws.application.delete,
        name='pbehavior/delete',
        payload=['_id']
    )
    def delete(_id):
        return rhpb.delete(_id)

    @route(
        ws.application.post,
        name='pbehavior/comment/create',
        payload=['pbehavior_id', 'author', 'message']
    )
    def create_comment(pbehavior_id, author, message):
        return rhpb.create_comment(pbehavior_id, author, message)

    @route(
        ws.application.put,
        name='pbehavior/comment/update',
        payload=['pbehavior_id', '_id', 'author', 'message']
    )
    def update_comment(pbehavior_id, _id, author=None, message=None):
        return rhpb.update_comment(pbehavior_id, _id, author, message)

    @route(
        ws.application.delete,
        name='pbehavior/comment/delete',
        payload=['pbehavior_id', '_id']
    )
    def delete_comment(pbehavior_id, _id):
        return rhpb.delete_comment(pbehavior_id, _id)
