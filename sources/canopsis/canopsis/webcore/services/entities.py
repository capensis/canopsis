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

from canopsis.common.ws import route


def get_entity(ws, mfilter):
    backend = ws.db.get_backend('entities')
    return backend.find_one(mfilter)


def get_entities(ws, mfilter, skip=0, limit=0):
    backend = ws.db.get_backend('entities')

    entities = backend.find(mfilter, skip=skip, limit=limit)

    total = entities.count()
    data = []

    for entity in entities:
        tmp = entity
        tmp['_id'] = str(tmp['_id'])

        data.append(tmp)

    return (data, total)


def exports(ws):
    @route(ws.application.get, payload=['start', 'limit'])
    def entities(etype=None, ename=None, start=0, limit=0):
        efilter = {}

        if etype is not None:
            efilter['type'] = etype

            if ename is not None:
                identifier = 'name'

                if etype == 'ack':
                    identifier = 'timestamp'

                efilter[identifier] = ename

        return get_entities(ws, efilter, skip=start, limit=limit)

    @route(ws.application.post, payload=['filter', 'start', 'limit'])
    def entities(filter, start=0, limit=0):
        return get_entities(ws, filter, skip=start, limit=limit)
