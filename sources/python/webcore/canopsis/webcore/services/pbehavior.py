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

from canopsis.common.ws import route
from canopsis.pbehavior.manager import PBehaviorManager


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
        return pbm.create(
            name=name, filter=filter, author=author,
            tstart=tstart, tstop=tstop, rrule=rrule,
            enabled=enabled, comments=comments,
            connector=connector, connector_name=connector_name
        )

    @route(
        ws.application.get,
        name='pbehavior/read',
        payload=['_id']
    )
    def read(_id=None):
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
        params.pop('_id')
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
        return pbm.create_pbehavior_comment(pbehavior_id, author, message)

    @route(
        ws.application.put,
        name='pbehavior/comment/update',
        payload=['pbehavior_id', '_id', 'author', 'message']
    )
    def update_comment(pbehavior_id, _id, author=None, message=None):
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

