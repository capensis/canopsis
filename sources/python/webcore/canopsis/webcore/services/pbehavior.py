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
            'enabled',
            'connector', 'connector_name'
        ]
    )
    def create(
            name, filter, author,
            tstart, tstop, rrule=None,
            enabled=True,
            connector='canopsis', connector_name='canopsis'
    ):
        return pbm.create(
            name=name, filter_=filter, author=author,
            tstart=tstart, tstop=tstop, rrule=rrule,
            enabled=enabled,
            connector=connector, connector_name=connector_name
        )

    @route(
        ws.application.get,
        name='pbehavior/read',
        payload=['_id']
    )
    def read(_id=None):
        return 'read'

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
            enabled=None
    ):
        return 'update'

    @route(
        ws.application.delete,
        name='pbehavior/delete',
        payload=['_id']
    )
    def delete(_id):
        return 'delete'

    @route(
        ws.application.post,
        name='pbehavior/comment/create',
        payload=['pbehavior_id', 'author', 'message']
    )
    def create_comment(pbehavior_id, author, message):
        return 'create comment'

    @route(
        ws.application.put,
        name='pbehavior/comment/update',
        payload=['pbehavior_id', '_id', 'auhtor', 'message']
    )
    def update_comment(pbehavior_id, _id, author=None, message=None):
        return 'update comment'

    @route(
        ws.application.delete,
        name='pbehavior/comment/delete',
        payload=['pbehavior_id', '_id']
    )
    def delete_comment(pbehavior_id, _id):
        return 'delete comment'
